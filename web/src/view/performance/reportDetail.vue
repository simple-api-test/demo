<template>
  <div style="background-color: #ffffff">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="任务名称：">
          <span>{{ reportName }}</span>
        </el-form-item>
        <el-form-item label="任务状态：">
          <a-tag :key="stateStyle" :color="stateStyle">{{ stateName }}</a-tag>
        </el-form-item>
        <el-form-item label="操作：">
          <el-button
            @click="updateDetail"
            type="primary"
            :disabled="boomerButton"
            >手动刷新</el-button
          >
          <el-button @click="stopBoomer" type="danger" :disabled="boomerButton"
            >停止运行</el-button
          >
        </el-form-item>

        <el-form-item>
          <el-select
            :disabled="boomerButton"
            v-model="timerValue"
            @change="updateTimerData(timerValue)"
            class="m-2"
            placeholder="Select"
            size="small"
          >
            <el-option
              v-for="item in timerOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <div style="padding: 10px">
      <el-tabs
        v-model="activeName"
        type="border-card"
        class="demo-tabs"
        @tab-click="handleClick"
      >
        <el-tab-pane label="压测指标" name="first">
          <div>
            <div style="padding-left: 10px; padding-right: 10px">
              <el-row :gutter="20" class="total">
                <el-col :span="6">
                  <div class="grid-content ep-bg-purple total-right">
                    <div>
                      <p>平均响应时间(ms)</p>
                    </div>
                    <div>
                      <h1>
                        {{ (respData.total_avg_response_time * 1).toFixed(2) }}
                      </h1>
                    </div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="grid-content ep-bg-purple total-right">
                    <div>
                      <p>响应时间PCT95(ms)</p>
                    </div>
                    <div>
                      <h1>{{ Math.ceil(PCT95) }}</h1>
                    </div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="grid-content ep-bg-purple total-right">
                    <div>
                      <p>平均RPS</p>
                    </div>
                    <div>
                      <h1>{{ (respData.total_rps * 1).toFixed(2) }}</h1>
                    </div>
                  </div>
                </el-col>
                <el-col :span="6">
                  <div class="grid-content ep-bg-purple">
                    <div>
                      <p>接口失败率(%)</p>
                    </div>
                    <div>
                      <h1>
                        {{ (respData.total_fail_ratio * 100).toFixed(2) }}
                      </h1>
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>

            <el-row :gutter="20">
              <el-col :span="24">
                <div id="sequential" class="grid-content ep-bg-purple"></div>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="24">
                <div id="transaction" class="grid-content ep-bg-purple"></div>
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>
        <el-tab-pane label="请求日志" name="second">
          <div style="padding-left: 40px; padding-right: 40px">
            <el-row :gutter="20">
              <el-col :span="24">
                <div class="grid-content ep-bg-purple">
                  <el-table :data="tableData" style="width: 100%">
                    <el-table-column
                      prop="method"
                      label="类型"
                      min-width="60"
                    />
                    <el-table-column
                      prop="name"
                      label="名称"
                      min-width="200"
                      max-width="400"
                    />
                    <el-table-column label="最小响应时间" min-width="120">
                      <template #default="scope">
                        {{ scope.row.min_response_time }} ms
                      </template>
                    </el-table-column>
                    <el-table-column label="最大响应时间" min-width="120">
                      <template #default="scope">
                        {{ scope.row.max_response_time }} ms
                      </template>
                    </el-table-column>
                    <el-table-column label="响应时间平均值" min-width="120">
                      <template #default="scope">
                        {{ scope.row.total_response_time }} ms
                      </template>
                    </el-table-column>
                    <el-table-column label="响应时间中位数" min-width="120">
                      <template #default="scope">
                        {{ scope.row.min_response_time }} ms
                      </template>
                    </el-table-column>
                    <el-table-column
                      prop="num_failures"
                      label="错误次数"
                      min-width="60"
                    />
                    <el-table-column label="错误率" min-width="60">
                      <template #default="scope">
                        {{
                          (
                            (scope.row.num_failures / scope.row.num_requests) *
                            100
                          ).toFixed(2)
                        }}
                        %
                      </template>
                    </el-table-column>
                    <el-table-column
                      prop="num_requests"
                      label="请求次数"
                      min-width="60"
                    />
                  </el-table>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>
        <el-tab-pane label="错误日志" name="third">
          <div style="padding-left: 40px; padding-right: 40px">
            <el-row :gutter="20">
              <el-col :span="24">
                <el-table :data="errData" style="width: 100%">
                  <el-table-column
                    prop="name"
                    label="名称"
                    min-width="200"
                    max-width="200"
                  />
                  <el-table-column
                    prop="occurrences"
                    label="错误次数"
                    width="200"
                  />
                  <el-table-column
                    prop="error"
                    label="错误信息"
                    min-width="600"
                  />
                </el-table>
              </el-col>
            </el-row>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
export default {
  name: "pReportDetail",
};
</script>

<script setup>
import { findReport } from "@/api/performance";
import { useRoute } from "vue-router";
import { onMounted, ref } from "vue";
import { Discount } from "@element-plus/icons-vue";
import * as echarts from "echarts/core";
import {
  ToolboxComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
} from "echarts/components";
import { LineChart } from "echarts/charts";
import { UniversalTransition } from "echarts/features";
import { CanvasRenderer } from "echarts/renderers";
import { runBoomer } from "@/api/runTestCase";
import { ElMessage } from "element-plus";

echarts.use([
  ToolboxComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  LineChart,
  CanvasRenderer,
  TitleComponent,
  UniversalTransition,
]);

const activeName = ref("first");
let performance_id = 0;
const handleClick = (Event) => {
  updateDetail();
};
const boomerButton = ref(false);
const route = useRoute();
const state = ref(0);
const respData = ref({});
const errData = ref([]);
const errDataKey = ref({});
const reportName = ref("");
const PCT95 = ref(0);
let reportID = 1;
const getTestCaseDetailFunc = async (testCaseID) => {
  const res = await findReport({ ID: testCaseID });
  if (res.code === 0) {
    performance_id = res.data.reapicase.performance_id;
    respData.value = res.data.reapicase;
    state.value = res.data.reapicase.state;
    runState(res.data.reapicase.state);
    reportName.value = res.data.reapicase.name;
    echartData(res.data.reapicase);
    transactionData(res.data.reapicase);
  }
};

const updateDetail = async () => {
  if (detailId === 0) {
    return;
  }
  const res = await findReport({ ID: reportID, DetailID: detailId });
  if (res.code === 0) {
    respData.value = res.data.reapicase;
    state.value = res.data.reapicase.state;
    runState(res.data.reapicase.state);
    echartData(res.data.reapicase);
    transactionData(res.data.reapicase);
  }
};

const updateUser = async () => {
  let data = {
    caseID: performance_id,
    run_type: 6,
    operation: {
      running: 2,
      spawnCount: 50,
      spawnRate: 5.0,
    },
  };
  const res = await runBoomer(data);
  if (res.code === 0) {
    ElMessage({
      type: "success",
      message: "运行成功",
    });
  }
};

const stopBoomer = async () => {
  let data = {
    caseID: performance_id,
    run_type: 6,
    operation: {
      running: 3,
    },
  };
  const res = await runBoomer(data);
  if (res.code === 0) {
    ElMessage({
      type: "success",
      message: "操作成功",
    });
  }
};

const stateName = ref("");
const stateStyle = ref("");
const runState = (t) => {
  if (t === 1) {
    stateName.value = "\xa0准备中\xa0";
    stateStyle.value = "green";
    return;
  }
  if (t === 2) {
    stateName.value = "\xa0运行中\xa0";
    stateStyle.value = "arcoblue";
    return;
  }
  if (t === 3) {
    stateName.value = "\xa0运行中\xa0";
    stateStyle.value = "arcoblue";
    return;
  }
  if (t === 4) {
    stateName.value = "\xa0停止中\xa0";
    stateStyle.value = "orange";
    return;
  }
  if (t === 5) {
    stateName.value = "\xa0已完成\xa0";
    stateStyle.value = "pinkpurple";
    boomerButton.value = true;
    return;
  }
  stateName.value = "未知状态";
};

const initData = async () => {
  if (route.params.id > 0) {
    reportID = route.params.id;
  }
  await getTestCaseDetailFunc(reportID);
};
initData();

let xTimeData = [];
let start_time = 0;
let xTimeNum = 0;
let maxUser = 0;
let maxRespTime = 0;
let maxFailRatio = 0;
let maxRPS = 0;

let dataUser = [];
let dataRespTime = [];
let dataFailRatio = [];
let dataRPS = [];

// 定时刷新
let timerData = setInterval(function () {}, 10000000000);
const timerOptions = [
  {
    value: 5,
    label: "5s",
  },
  {
    value: 15,
    label: "15s",
  },
  {
    value: 30,
    label: "30s",
  },
  {
    value: 30,
    label: "60s",
  },
  {
    value: 0,
    label: "关闭自动刷新",
  },
];
const timerValue = ref("手动刷新");
const updateTimerData = (timer) => {
  clearInterval(timerData);
  if (timer === 0 || state.value === 5) {
    return;
  }
  timerData = setInterval(function () {
    if (state.value === 5) {
      clearInterval(timerData);
    }
    updateDetail();
  }, timer * 1000);
};

let colorUser = "#FF9A2E";
let colorSequential = "#00B42A";
let colorArea = "#E8FFEA";

let detailId = 0;
const tableData = ref([]);

const echartData = (data) => {
  // PCT95.value = quantile(pctList, .95)

  data.performance_report_detail.forEach((item, index) => {
    PCT95.value = item.stats_total.pct["95"];
    tableData.value = [];
    let err = item.errors;
    for (let key in err) {
      if (err[key].method === "request") {
        if (errDataKey.value[key]) {
          errData.value[errDataKey.value[key] - 1].occurrences +=
            err[key].occurrences;
        } else {
          errDataKey.value[key] = errData.value.length + 1;
          let data = {
            name: err[key].name,
            error: err[key].error,
            occurrences: err[key].occurrences,
            method: err[key].method,
          };
          errData.value.push(data);
        }
      }
    }
    item.stats.forEach((item2, index2) => {
      if (
        item2.method === "testcase" ||
        (item2.method === "transaction" && item2.name === "Action")
      ) {
      } else {
        tableData.value.push(item2);
      }
    });

    detailId = item.ID;
    if (start_time === 0) {
      start_time = item.stats_total.start_time;
    }
    if (item.user_count > maxUser) {
      maxUser = getMax(item.user_count);
    }
    dataUser.push(item.user_count);
    if (item.total_avg_response_time > maxRespTime) {
      maxRespTime = getMax(item.total_avg_response_time);
    }
    if (item.total_fail_ratio * 100 > maxFailRatio) {
      maxFailRatio = getMax(item.total_fail_ratio * 100);
    }
    dataFailRatio.push(Number((item.total_fail_ratio * 100).toFixed(2)));
    xTime(xTimeNum);
    dataRespTime.push(Number(item.total_avg_response_time.toFixed(2)));
    if (item.total_rps > maxRPS) {
      maxRPS = getMax(Math.ceil(item.total_rps));
    }
    dataRPS.push(Number(item.total_rps.toFixed(2)));

    xTimeNum++;
  });

  sequentialEchart();
};

let transactionSeriesKey = { Action: 0 };
let transactionName = ["所有事务"];
let transactionSeriesID = 0;
let transactionExecuted = [[]];
let transactionSuccess = [[]];
let transactionSuccessRatio = [[]];
let transactionTPS = [[]];
let transactionFail = [[]];
let maxTransaction = 0;
let maxTransactionSuccess = 0;
let maxTPS = 0;

const deepCopy = (obj) => {
  return JSON.parse(JSON.stringify(obj));
};

const transactionData = (data) => {
  let xAxis = [];
  let xAxi = {
    gridIndex: 0,
    type: "category",
    scale: true,
    axisTick: {
      alignWithLabel: true,
    },
    splitLine: {
      show: false,
    },
    data: xTimeData,
  };
  let series = [];
  let serieUser = {
    name: "并发用户数",
    symbol: "none",
    type: "line",
    xAxisIndex: 0,
    yAxisIndex: 1,
    hoverAnimation: true,
    data: dataUser,
    itemStyle: {
      normal: {
        lineStyle: {
          width: 2,
          type: "dashed",
        },
        color: colorUser,
      },
    },
  };

  for (let i = 0; i < 4; i++) {
    xAxi.gridIndex = i;
    xAxis.push(deepCopy(xAxi));
  }

  data.performance_report_detail.forEach((item, index) => {
    let i = transactionExecuted[0].length;
    let actionExecuted = 0;
    let actionSuccess = 0;
    item.stats.forEach((item2, index2) => {
      if (item2.method === "transaction" && item2.name !== "Action") {
        if (transactionSeriesKey[item2.name] === undefined) {
          transactionName.push(item2.name);
          transactionSeriesID++;
          transactionSeriesKey[item2.name] = transactionSeriesID;
          transactionExecuted.push([]);
          transactionSuccess.push([]);
          transactionSuccessRatio.push([]);
          transactionTPS.push([]);
        }
        // 执行事务数
        let transactionExecutedLength =
          transactionExecuted[transactionSeriesKey[item2.name]].length;
        if (transactionExecutedLength > 0) {
          let num =
            item2.num_requests +
            transactionExecuted[transactionSeriesKey[item2.name]][
              transactionExecutedLength - 1
            ];
          actionExecuted += num;
          transactionExecuted[transactionSeriesKey[item2.name]][i] = num;
        } else {
          transactionExecuted[transactionSeriesKey[item2.name]][i] =
            item2.num_requests;
          actionExecuted += item2.num_requests;
        }

        // 成功数、成功率
        let transactionSuccessLength =
          transactionSuccess[transactionSeriesKey[item2.name]].length;
        let successNum = item2.num_requests - item2.num_failures;
        let successRatio = (
          (successNum / (item2.num_requests * 1.0)) *
          100
        ).toFixed(2);
        transactionSuccessRatio[transactionSeriesKey[item2.name]][i] =
          successRatio;
        if (transactionSuccessLength > 0) {
          let num =
            successNum +
            transactionSuccess[transactionSeriesKey[item2.name]][
              transactionExecutedLength - 1
            ];
          transactionSuccess[transactionSeriesKey[item2.name]][i] = num;
          actionSuccess += num;
        } else {
          transactionSuccess[transactionSeriesKey[item2.name]][i] = successNum;
          actionSuccess += successNum;
        }

        // TPS
        if (item2.current_rps > maxTPS) {
          maxTPS = getMax(item2.current_rps);
        }
        transactionTPS[transactionSeriesKey[item2.name]][i] = Number(
          item2.current_rps.toFixed(2)
        );
        // if (item2.current_rps.toFixed){
        //   transactionTPS[transactionSeriesKey[item2.name]][i]=Number(item2.current_rps.toFixed(2))
        // }else {
        //
        //   transactionTPS[transactionSeriesKey[item2.name]][i]=item2.current_rps
        // }
      }

      if (index2 === item.stats.length - 1) {
        if (actionExecuted > maxTransaction) {
          maxTransaction = getMax(actionExecuted);
        }
        transactionExecuted[0][i] = actionExecuted;

        if (actionSuccess > maxTransactionSuccess) {
          maxTransactionSuccess = getMax(actionSuccess);
        }
        transactionSuccess[0][i] = actionSuccess;

        let successNum = item2.num_requests - item2.num_failures;
        let successRatio = (
          (successNum / (item2.num_requests * 1.0)) *
          100
        ).toFixed(2);
        transactionSuccessRatio[0][i] = successRatio;

        transactionTPS[0][i] = Number(item2.current_rps.toFixed(2));
      }
    });
  });

  let serie = {
    name: "响应时间",
    symbol: "none",
    type: "line",
    xAxisIndex: 0,
    yAxisIndex: 0,
    hoverAnimation: true,
    data: dataRespTime,
  };

  // 执行事务数折线图数据
  transactionExecuted.forEach((item, index) => {
    serie.data = item;
    serie.name = transactionName[index];
    series.push(JSON.parse(JSON.stringify(serie)));
  });
  serieUser.xAxisIndex = 0;
  serieUser.yAxisIndex = 1;
  series.push(deepCopy(serieUser));

  // 事务成功数折线图数据
  transactionSuccess.forEach((item, index) => {
    serie.data = item;
    serie.xAxisIndex = 1;
    serie.yAxisIndex = 2;
    serie.name = transactionName[index];
    series.push(JSON.parse(JSON.stringify(serie)));
  });
  serieUser.xAxisIndex = 1;
  serieUser.yAxisIndex = 3;
  series.push(deepCopy(serieUser));

  // 事务成功率折线图数据
  transactionSuccessRatio.forEach((item, index) => {
    serie.data = item;
    serie.xAxisIndex = 2;
    serie.yAxisIndex = 4;
    serie.name = transactionName[index];
    series.push(JSON.parse(JSON.stringify(serie)));
  });
  serieUser.xAxisIndex = 2;
  serieUser.yAxisIndex = 5;
  series.push(deepCopy(serieUser));

  // TPS 折线图数据
  transactionTPS.forEach((item, index) => {
    serie.data = item;
    serie.xAxisIndex = 3;
    serie.yAxisIndex = 6;
    serie.name = transactionName[index];
    series.push(JSON.parse(JSON.stringify(serie)));
  });
  serieUser.xAxisIndex = 3;
  serieUser.yAxisIndex = 7;
  series.push(deepCopy(serieUser));

  transactionName.push("并发用户数");
  transactionEchart(xAxis, series, transactionName);
};

const xTime = (i) => {
  let runTime = i * 3;
  let hour = Math.floor(runTime / 3600);
  runTime = runTime % 3600;
  let minute = Math.floor(runTime / 60);
  runTime = runTime % 60;
  let second = runTime;
  xTimeData.push(`${hour}:${minute}:${second}`);
};

const getMax = (max) => {
  max = Math.ceil(max);
  var bite = 0;
  if (max < 1) {
    return 1.0;
  }
  if (max < 5) {
    return 5;
  }
  if (max < 10) {
    return 10;
  }
  if (max >= 2500) {
    while (max >= 100) {
      max /= 100;
      bite += 1;
    }
    return Math.ceil(max) * Math.pow(100, bite);
  }
  if (max >= 500) {
    while (max >= 50) {
      max /= 50;
      bite += 1;
    }
    return Math.ceil(max) * Math.pow(50, bite);
  }
  while (max >= 10) {
    max /= 10;
    bite += 1;
  }
  return Math.ceil(max) * Math.pow(10, bite);
};

const sequentialEchart = () => {
  const optionSequential = {
    grid: [
      {
        left: "5%",
        right: "69%",
        top: "20%",
        bottom: "20%",
      },
      {
        left: "37%",
        right: "37%",
        top: "20%",
        bottom: "20%",
      },
      {
        left: "69%",
        right: "5%",
        top: "20%",
        bottom: "20%",
      },
    ],
    tooltip: {
      trigger: "axis",
      backgroundColor: "#ccc",
      borderWidth: 1,
      borderColor: "#ccc",
      padding: 10,
      textStyle: {
        color: "#000",
        fontSize: 12,
        align: "left",
      },
    },
    xAxis: [
      {
        type: "category",
        scale: true,
        axisTick: {
          alignWithLabel: true,
        },
        splitLine: {
          show: false,
        },
        data: xTimeData,
      },
      {
        gridIndex: 1,
        type: "category",
        scale: true,
        axisTick: {
          alignWithLabel: true,
        },
        splitLine: {
          show: false,
        },
        data: xTimeData,
      },
      {
        gridIndex: 2,
        type: "category",
        scale: true,
        axisTick: {
          alignWithLabel: true,
        },
        splitLine: {
          show: false,
        },
        data: xTimeData,
      },
    ],
    yAxis: [
      {
        type: "value",
        name: "响应时间",
        scale: true,
        axisLabel: {
          formatter: "{value} ms",
        },
        min: 0,
        max: maxRespTime,
        interval: maxRespTime / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
      {
        type: "value",
        name: "业务场景RPS",
        gridIndex: 1,
        scale: true,
        splitNumber: 3,
        min: 0,
        max: maxRPS,
        interval: maxRPS / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        gridIndex: 1,
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
      {
        type: "value",
        name: "接口错误率",
        gridIndex: 2,
        scale: true,
        splitNumber: 3,
        axisLabel: {
          formatter: "{value} %",
        },
        min: 0,
        max: maxFailRatio,
        interval: maxFailRatio / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        gridIndex: 2,
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
    ],
    dataZoom: [],
    series: [
      {
        name: "响应时间",
        itemStyle: { color: colorSequential },
        areaStyle: { color: colorArea },
        symbol: "none",
        type: "line",
        xAxisIndex: 0,
        yAxisIndex: 0,
        hoverAnimation: true,
        data: dataRespTime,
      },
      {
        name: "并发用户数",
        symbol: "none",
        type: "line",
        xAxisIndex: 0,
        yAxisIndex: 1,
        hoverAnimation: true,
        data: dataUser,
        itemStyle: {
          normal: {
            lineStyle: {
              width: 2,
              type: "dashed",
            },
            color: colorUser,
          },
        },
      },
      {
        name: "RPS",
        itemStyle: { color: colorSequential },
        areaStyle: { color: colorArea },
        symbol: "none",
        type: "line",
        xAxisIndex: 1,
        yAxisIndex: 2,
        hoverAnimation: true,
        data: dataRPS,
      },
      {
        name: "并发用户数",
        symbol: "none",
        type: "line",
        xAxisIndex: 1,
        yAxisIndex: 3,
        hoverAnimation: true,
        data: dataUser,
        itemStyle: {
          normal: {
            lineStyle: {
              width: 2,
              type: "dashed",
            },
            color: colorUser,
          },
        },
      },
      {
        name: "接口报错率",
        itemStyle: { color: colorSequential },
        areaStyle: { color: colorArea },
        symbol: "none",
        type: "line",
        xAxisIndex: 2,
        yAxisIndex: 4,
        hoverAnimation: true,
        data: dataFailRatio,
      },
      {
        name: "并发用户数",
        symbol: "none",
        type: "line",
        xAxisIndex: 2,
        yAxisIndex: 5,
        hoverAnimation: true,
        data: dataUser,
        itemStyle: {
          normal: {
            lineStyle: {
              width: 2,
              type: "dashed",
            },
            color: colorUser,
          },
        },
      },
    ],
  };
  let sequentialChart = document.getElementById("sequential");
  let myChart = echarts.init(sequentialChart);
  optionSequential && myChart.setOption(optionSequential);
};

const transactionEchart = (xAxis, series, legend) => {
  const transaction = {
    title: [
      {
        text: "执行事务数",
        left: "20%",
      },
      {
        text: "事务成功数",
        right: "20%",
      },
      {
        text: "事务成功率",
        top: "52%",
        left: "20%",
      },
      {
        text: "TPS",
        right: "20%",
        top: "52%",
      },
    ],
    legend: {
      data: legend,
    },
    grid: [
      {
        left: "4%",
        right: "54%",
        top: "10%",
        bottom: "55%",
      },
      {
        left: "54%",
        right: "4%",
        top: "10%",
        bottom: "55%",
      },
      {
        left: "4%",
        right: "54%",
        top: "62%",
        bottom: "3%",
      },
      {
        left: "54%",
        right: "4%",
        top: "62%",
        bottom: "3%",
      },
    ],
    tooltip: {
      trigger: "axis",
      backgroundColor: "#ccc",
      borderWidth: 1,
      borderColor: "#ccc",
      padding: 10,
      textStyle: {
        color: "#000",
        fontSize: 12,
        align: "left",
      },
    },
    xAxis: xAxis,
    yAxis: [
      {
        type: "value",
        name: "执行事务数",
        scale: true,
        axisLabel: {
          formatter: "{value}",
        },
        min: 0,
        max: maxTransaction,
        interval: maxTransaction / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
      {
        type: "value",
        name: "事务成功数",
        gridIndex: 1,
        scale: true,
        splitNumber: 3,
        min: 0,
        max: maxTransactionSuccess,
        interval: maxTransactionSuccess / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        gridIndex: 1,
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
      {
        type: "value",
        name: "事务成功率",
        gridIndex: 2,
        scale: true,
        splitNumber: 3,
        axisLabel: {
          formatter: "{value} %",
        },
        min: 0,
        max: 100,
        interval: 100 / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        gridIndex: 2,
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
      {
        type: "value",
        name: "TPS",
        gridIndex: 3,
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxTPS,
        interval: maxTPS / 5,
      },
      {
        type: "value",
        name: "并发用户数",
        gridIndex: 3,
        scale: true,
        splitNumber: 3,
        splitLine: {
          lineStyle: {
            type: "dashed",
          },
          show: true,
        },
        min: 0,
        max: maxUser,
        interval: maxUser / 5,
      },
    ],
    dataZoom: [],
    series: series,
  };
  let transactionChart = document.getElementById("transaction");
  let myChart = echarts.init(transactionChart);
  transaction && myChart.setOption(transaction);
};
</script>

<style scoped>
.el-row {
  margin-bottom: 20px;
}
.el-row:last-child {
  margin-bottom: 0;
}
.el-col {
  border-radius: 4px;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}
.total {
  border-style: solid;
  border-width: thin;
  border-color: #cccccc;
  border-radius: 5px;
}

.total-right {
  border-right: 1px solid #cccccc;
  border-radius: 0px;
}

h1 {
  text-align: center;
}

p {
  margin-top: 25px;
  margin-bottom: 5px;
  text-align: center;
  font-size: 16px;
}

#response_time,
#fail_ratio,
#rps {
  height: 300px;
}

#sequential {
  height: 300px;
  border-style: solid;
  border-width: 1px;
  border-color: #c9cdd4;
}

#transaction {
  height: 650px;
  border-style: solid;
  border-width: 1px;
  border-color: #c9cdd4;
}
</style>
