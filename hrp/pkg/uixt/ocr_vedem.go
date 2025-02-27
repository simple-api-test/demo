package uixt

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/test-instructor/yangfan/hrp/internal/builtin"
	"github.com/test-instructor/yangfan/hrp/internal/code"
	"github.com/test-instructor/yangfan/hrp/internal/env"
	"github.com/test-instructor/yangfan/hrp/internal/json"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

type OCRResult struct {
	Text   string   `json:"text"`
	Points []PointF `json:"points"`
}

type ResponseOCR struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	OCRResult []OCRResult `json:"ocrResult"`
}

type veDEMOCRService struct{}

func newVEDEMOCRService() (*veDEMOCRService, error) {
	if err := checkEnv(); err != nil {
		return nil, err
	}
	return &veDEMOCRService{}, nil
}

func checkEnv() error {
	if env.VEDEM_OCR_URL == "" {
		return errors.Wrap(code.OCREnvMissedError, "VEDEM_OCR_URL missed")
	}
	if env.VEDEM_OCR_AK == "" {
		return errors.Wrap(code.OCREnvMissedError, "VEDEM_OCR_AK missed")
	}
	if env.VEDEM_OCR_SK == "" {
		return errors.Wrap(code.OCREnvMissedError, "VEDEM_OCR_SK missed")
	}
	return nil
}

func (s *veDEMOCRService) getOCRResult(imageBuf []byte) ([]OCRResult, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("withDet", "true")
	// bodyWriter.WriteField("timestampOnly", "true")

	formWriter, err := bodyWriter.CreateFormFile("image", "screenshot.png")
	if err != nil {
		return nil, errors.Wrap(code.OCRRequestError,
			fmt.Sprintf("create form file error: %v", err))
	}
	size, err := formWriter.Write(imageBuf)
	if err != nil {
		return nil, errors.Wrap(code.OCRRequestError,
			fmt.Sprintf("write form error: %v", err))
	}

	err = bodyWriter.Close()
	if err != nil {
		return nil, errors.Wrap(code.OCRRequestError,
			fmt.Sprintf("close body writer error: %v", err))
	}

	req, err := http.NewRequest("POST", env.VEDEM_OCR_URL, bodyBuf)
	if err != nil {
		return nil, errors.Wrap(code.OCRRequestError,
			fmt.Sprintf("construct request error: %v", err))
	}

	token := builtin.Sign("auth-v2", env.VEDEM_OCR_AK, env.VEDEM_OCR_SK, bodyBuf.Bytes())
	req.Header.Add("Agw-Auth", token)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())

	var resp *http.Response
	// retry 3 times
	for i := 1; i <= 3; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		var logID string
		if resp != nil {
			logID = getLogID(resp.Header)
		}
		log.Error().Err(err).
			Str("logID", logID).
			Int("imageBufSize", size).
			Msgf("request OCR service failed, retry %d", i)
		time.Sleep(1 * time.Second)
	}
	if resp == nil {
		return nil, code.OCRServiceConnectionError
	}

	defer resp.Body.Close()

	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(code.OCRResponseError,
			fmt.Sprintf("read response body error: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(code.OCRResponseError,
			fmt.Sprintf("unexpected response status code: %d, results: %v",
				resp.StatusCode, string(results)))
	}

	var ocrResult ResponseOCR
	err = json.Unmarshal(results, &ocrResult)
	if err != nil {
		return nil, errors.Wrap(code.OCRResponseError,
			fmt.Sprintf("json unmarshal response body error: %v", err))
	}

	return ocrResult.OCRResult, nil
}

func getLogID(header http.Header) string {
	if len(header) == 0 {
		return ""
	}

	logID, ok := header["X-Tt-Logid"]
	if !ok || len(logID) == 0 {
		return ""
	}
	return logID[0]
}

func (s *veDEMOCRService) FindText(text string, imageBuf []byte, options ...DataOption) (rect image.Rectangle, err error) {
	data := NewData(map[string]interface{}{}, options...)

	ocrResults, err := s.getOCRResult(imageBuf)
	if err != nil {
		log.Error().Err(err).Msg("getOCRResult failed")
		return
	}

	var rects []image.Rectangle
	var ocrTexts []string
	for _, ocrResult := range ocrResults {
		rect = image.Rectangle{
			// ocrResult.Points 顺序：左上 -> 右上 -> 右下 -> 左下
			Min: image.Point{
				X: int(ocrResult.Points[0].X),
				Y: int(ocrResult.Points[0].Y),
			},
			Max: image.Point{
				X: int(ocrResult.Points[2].X),
				Y: int(ocrResult.Points[2].Y),
			},
		}
		if rect.Min.X >= data.Scope[0] && rect.Max.X <= data.Scope[2] && rect.Min.Y >= data.Scope[1] && rect.Max.Y <= data.Scope[3] {
			ocrTexts = append(ocrTexts, ocrResult.Text)

			// not contains text
			if !strings.Contains(ocrResult.Text, text) {
				continue
			}

			rects = append(rects, rect)
		}

		// contains text while not match exactly
		if ocrResult.Text != text {
			continue
		}

		// match exactly, and not specify index, return the first one
		if data.Index == 0 {
			return rect, nil
		}
	}

	if len(rects) == 0 {
		return image.Rectangle{}, errors.Wrap(code.OCRTextNotFoundError,
			fmt.Sprintf("text %s not found in %v", text, ocrTexts))
	}

	// get index
	idx := data.Index
	if idx > 0 {
		// NOTICE: index start from 1
		idx = idx - 1
	} else if idx < 0 {
		idx = len(rects) + idx
	}

	// index out of range
	if idx >= len(rects) {
		return image.Rectangle{}, errors.Wrap(code.OCRTextNotFoundError,
			fmt.Sprintf("text %s found %d, index %d out of range", text, len(rects), idx))
	}

	return rects[idx], nil
}

func (s *veDEMOCRService) FindTexts(texts []string, imageBuf []byte, options ...DataOption) (rects []image.Rectangle, err error) {
	ocrResults, err := s.getOCRResult(imageBuf)
	if err != nil {
		log.Error().Err(err).Msg("getOCRResult failed")
		return
	}

	data := NewData(map[string]interface{}{}, options...)
	ocrTexts := map[string]bool{}

	var success bool
	var rect image.Rectangle
	for _, text := range texts {
		var found bool
		for _, ocrResult := range ocrResults {
			rect = image.Rectangle{
				// ocrResult.Points 顺序：左上 -> 右上 -> 右下 -> 左下
				Min: image.Point{
					X: int(ocrResult.Points[0].X),
					Y: int(ocrResult.Points[0].Y),
				},
				Max: image.Point{
					X: int(ocrResult.Points[2].X),
					Y: int(ocrResult.Points[2].Y),
				},
			}

			if rect.Min.X >= data.Scope[0] && rect.Max.X <= data.Scope[2] && rect.Min.Y >= data.Scope[1] && rect.Max.Y <= data.Scope[3] {
				ocrTexts[ocrResult.Text] = true

				// not contains text
				if !strings.Contains(ocrResult.Text, text) {
					continue
				}

				found = true
				rects = append(rects, rect)
				break
			}
		}
		if !found {
			rects = append(rects, image.Rectangle{})
		}
		success = found || success
	}

	if !success {
		return rects, errors.Wrap(code.OCRTextNotFoundError,
			fmt.Sprintf("texts %s not found in %v", texts, ocrTexts))
	}

	return rects, nil
}

type OCRService interface {
	FindText(text string, imageBuf []byte, index ...int) (rect image.Rectangle, err error)
}

func (dExt *DriverExt) FindTextByOCR(ocrText string, options ...DataOption) (x, y, width, height float64, err error) {
	var bufSource *bytes.Buffer
	if bufSource, err = dExt.takeScreenShot(); err != nil {
		err = fmt.Errorf("takeScreenShot error: %v", err)
		return
	}

	service, err := newVEDEMOCRService()
	if err != nil {
		return
	}
	rect, err := service.FindText(ocrText, bufSource.Bytes(), options...)
	if err != nil {
		log.Warn().Msgf("FindText failed: %s", err.Error())
		return
	}

	log.Info().Str("ocrText", ocrText).
		Interface("rect", rect).Msgf("FindTextByOCR success")
	x, y, width, height = dExt.MappingToRectInUIKit(rect)
	return
}

func (dExt *DriverExt) FindTextsByOCR(ocrTexts []string, options ...DataOption) (points [][]float64, err error) {
	var bufSource *bytes.Buffer
	if bufSource, err = dExt.takeScreenShot(); err != nil {
		err = fmt.Errorf("takeScreenShot error: %v", err)
		return
	}

	service, err := newVEDEMOCRService()
	if err != nil {
		return
	}
	rects, err := service.FindTexts(ocrTexts, bufSource.Bytes(), options...)
	if err != nil {
		log.Warn().Msgf("FindTexts failed: %s", err.Error())
		return
	}

	log.Info().Interface("ocrTexts", ocrTexts).
		Interface("rects", rects).Msgf("FindTextsByOCR success")
	for _, rect := range rects {
		x, y, width, height := dExt.MappingToRectInUIKit(rect)
		points = append(points, []float64{x, y, width, height})
	}

	return
}
