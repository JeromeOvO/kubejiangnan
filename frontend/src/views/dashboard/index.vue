<template>
  <div class="dashboard-container">
    <div class="dashboard-group">
      <el-row>
        <el-col class="card" v-for="item in cluster" :span="6">
          <dashboard-card :data="item"/>
        </el-col>
      </el-row>
    </div>
    <div class="dashboard-group">
      <el-row>
        <el-col class="card" :span="4" v-for="item in resource">
          <dash-resource-card :data="item"/>
        </el-col>
      </el-row>
    </div>
    <div class="dashboard-group">
      <el-row style="height: 200px;">
        <el-col v-for="item in usage" class="card" :span="8">
          <dash-pointer-card :data="item"/>
        </el-col>
      </el-row>
    </div>
    <div>
      <el-row style="height: 200px;">
        <el-col v-for="item in usageRange" class="card" :span="12">
          <dash-range-card :data="item"/>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import {mapGetters} from 'vuex'
import DashboardCard from "./components/DashboardCard";
import DashResourceCard from "./components/DashK8SResourceCard";
import DashPointerCard from "./components/DashPointerCard";
import DashRangeCard from "./components/DashRangeCard";

const randomRgb = () => {
  let R = Math.floor(Math.random() * 255);
  let G = Math.floor(Math.random() * 255);
  let B = Math.floor(Math.random() * 255);
  return R + ',' + G + ',' + B
}
const clusterData = [
  {
    name: "K8S",
    desc: "Cluster",
    color: "#409EFF"
  },
  {
    name: "v1.23.9",
    desc: "Kubernetes Version",
    color: "#409EFF"
  },
  {
    name: "3",
    desc: "Nodes",
    color: "#409EFF"
  },
  {
    name: "1个月前",
    desc: "Created",
    color: "#409EFF"
  }
]
const resourceData = [
  {
    name: "11",
    desc: "Namespaces",
    logo: "k8s"
  },
  {
    name: "35",
    desc: "Pods",
    logo: "pod"
  },
  {
    name: "2",
    desc: "Ingresses",
    logo: "svc"
  },
  {
    name: "6",
    desc: "Services",
    logo: "svc"
  },
  {
    name: "4",
    desc: "PVC",
    logo: "data",
  },
  {
    name: "20",
    desc: "ConfigMaps",
    logo: "data"
  },
  {
    name: "21",
    desc: "Secrets",
    logo: "data"
  },
  {
    name: "8",
    desc: "Deployments",
    logo: "workload"
  },
  {
    name: "1",
    desc: "StatefulSets",
    logo: "workload"
  },
  {
    name: "1",
    desc: "CronJobs",
    logo: "workload"
  },
  {
    name: "2",
    desc: "DaemonSets",
    logo: "workload"
  },
  {
    name: "6",
    desc: "CRDS",
    logo: "crds"
  }
]
const usagePointerDemo = [{
  title: "Pod使用量",
  value: 32,
  // height: "200px",
  bgc: "#F2F6FC",
  // dc: "red"
},
  {
    title: "内存使用量",
    value: 79,
    // height: "200px",
    bgc: "#F2F6FC",
    // dc: "red"
  },
  {
    title: "CPU使用量",
    value: 84,
    // height: "200px",
    bgc: "#F2F6FC",
    // dc: "red"
  }]
const usageRangeDemo = [
  {
    title: "CPU变化趋势",
    height: "240px",
    x: ["10:11", "10:12", "10:13", "10:14", "10:15", "10:16", "10:17", "10:18", "10:19", "10:20", "10:11", "10:12", "10:13", "10:14", "10:15", "10:16", "10:17", "10:18", "10:19", "10:20"],
    y: [300, 400, 399, 400, 500, 300, 400, 399, 400, 500, 300, 400, 399, 400, 500, 800, 900, 1000, 999, 999],
    bgc: "#F2F6FC",
    dc: "red"
  },
  {
    title: "内存变化趋势",
    height: "240px",
    x: ["10:11", "10:12", "10:13", "10:14", "10:15", "10:16", "10:17", "10:18", "10:19", "10:20", "10:11", "10:12", "10:13", "10:14", "10:15", "10:16", "10:17", "10:18", "10:19", "10:20"],
    y: [1024, 1024, 1024, 1024, 1024, 1024, 400, 399, 400, 500, 300, 400, 399, 400, 500, 800, 900, 1000, 999, 999],
    bgc: "#F2F6FC",
    dc: "red"
  }
]
for (let i = resourceData.length - 1; i >= 0; i--) {
  resourceData[i].color = randomRgb()
}

for (let i = clusterData.length - 1; i >= 0; i--) {
  clusterData[i].color = randomRgb()
}
export default {
  name: 'Dashboard',
  components: {
    DashboardCard,
    DashResourceCard,
    DashPointerCard,
    DashRangeCard
  },
  data() {
    return {
      resource: [],
      cluster: [],
      usage: [],
      usageRange: []
    }
  },
  created() {
    this.getData()
  },
  computed: {
    ...
      mapGetters([
        'name'
      ])
  },
  methods: {
    getData() {
      this.$store.dispatch("metrics/getMetricsDashboard").then(res => {
        this.cluster = res.data.cluster
        this.resource = res.data.resource
        this.usage = res.data.usage
        this.usageRange = res.data.usageRange
      })
    }
  }
}
</script>

<style lang="scss" scoped>

.card {
  padding: 5px;
}

.dashboard {
  &-container {
    margin: 30px;
  }

  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
