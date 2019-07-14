<template>
  <el-row class="app-container">
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="160" />
      <el-table-column prop="name" label="项目名称" />
      <el-table-column prop="publisherName" label="构建者" width="160" />
      <el-table-column prop="updateTime" label="上次构建时间" width="160" />
      <el-table-column prop="operation" label="操作" width="260">
        <template slot-scope="scope">
          <el-button size="small" type="primary" @click="publish(scope.row.id)">构建</el-button>
          <el-button size="small" type="success" @click="handleDetail(scope.row.id)">详情</el-button>
          <el-button size="small" type="danger">回滚</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="上一次构建记录" :visible.sync="dialogVisible">
      <el-row class="project-detail">
        <el-row>
          <el-row style="margin:5px 0">git同步信息</el-row>
          <el-row style="margin:5px 0">状态: {{ gitTrace['state'] === 1 ? '成功' : '失败' }}</el-row>
          <el-row style="margin:5px 0" v-html="formatDetail(gitTrace['detail'])" />
        </el-row>
        <hr>
        <el-row>
          <el-row style="margin:5px 0">rsync同步服务器信息</el-row>
          <el-row v-for="(item, index) in rsyncTraceList" :key="index">
            <el-row style="margin:5px 0">服务器: {{ item['serverName'] }}</el-row>
            <el-row style="margin:5px 0">状态: {{ item['state'] === 1 ? '成功' : '失败' }}</el-row>
            <el-row style="margin:5px 0" v-html="formatDetail(item['detail'])" />
            <hr>
          </el-row>
        </el-row>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { get, getDetail, publish } from '@/api/deploy'
import { parseTime } from '@/utils'

const STATE = ['构建中', '构建成功', '构建失败', '撤回']
export default {
  data() {
    return {
      dialogVisible: false,
      tableData: [],
      gitTrace: {},
      rsyncTraceList: []

    }
  },
  created() {
    this.get()
  },
  methods: {
    get() {
      get().then((response) => {
        const projectList = response.data.projectList
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
          element.state = STATE[element.state]
        })
        this.tableData = projectList
      })
    },
    publish(id) {
      publish(id).then((response) => {
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
        this.get()
      })
    },
    handleDetail(id) {
      getDetail(id).then((response) => {
        this.dialogVisible = true
        this.gitTrace = response.data.gitTrace
        this.rsyncTraceList = response.data.rsyncTraceList
      })
    },
    formatDetail(detail) {
      return detail ? detail.replace(/\n/g, '<br>') : ''
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import "@/styles/mixin.scss";
.project-detail {
  height:300px;
  overflow-y: auto;
  @include scrollBar();
}
</style>