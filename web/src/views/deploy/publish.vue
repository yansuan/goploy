<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex">
      <el-select v-model="groupId" placeholder="选择分组" @change="handleGroupChange">
        <el-option label="默认" :value="0" />
        <el-option
          v-for="(item, index) in groupOption"
          :key="index"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
      <el-input v-model="projectName" style="width:300px" placeholder="请输入项目名称" @change="getList" />
    </el-row>
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%;margin-top: 5px;"
    >
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="项目名称" min-width="160" show-overflow-tooltip>
        <template slot-scope="scope">
          <b v-if="scope.row.environment === '生产环境'" style="color: #F56C6C">{{ scope.row.name }} - {{ scope.row.environment }}</b>
          <b v-else-if="scope.row.environment === '测试环境'" style="color: #E6A23C">{{ scope.row.name }} - {{ scope.row.environment }}</b>
          <b v-else style="color: #909399">{{ scope.row.name }} - {{ scope.row.environment }}</b>
        </template>
      </el-table-column>
      <el-table-column prop="group" label="分组">
        <template slot-scope="scope">
          {{ findGroupName(scope.row.groupId) }}
        </template>
      </el-table-column>
      <el-table-column prop="branch" label="分支" />
      <el-table-column prop="publisherName" label="构建者" width="160" />
      <el-table-column prop="deployState" label="构建状态" width="70">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.deployState === 0" type="info" effect="plain">未构建</el-tag>
          <el-tag v-else-if="scope.row.deployState === 1" type="warning" effect="plain">构建中</el-tag>
          <el-tag v-else-if="scope.row.deployState === 2" type="success" effect="plain">成功</el-tag>
          <el-tag v-else type="danger" effect="plain">失败</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updateTime" label="上次构建时间" width="160" />
      <el-table-column prop="operation" label="操作" width="175" fixed="right">
        <template slot-scope="scope">
          <el-dropdown
            class="publish-btn"
            split-button
            trigger="click"
            :disabled="scope.row.deployState === 1"
            type="primary"
            @click="publish(scope.row)"
            @command="handlePublishCommand"
          >
            构建
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item :command="scope.row">选择具体commit</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>

          <el-button type="success" @click="handleDetail(scope.row)">详情</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="构建记录" :visible.sync="dialogVisible" class="publish-record">
      <el-row>
        <el-col :span="8">
          <el-radio-group v-model="publishToken" @change="handleDetailChange">
            <el-row v-for="(item, index) in gitTraceList" :key="index">
              <el-row style="margin:5px 0">
                <el-radio style="margin-left: 10px;margin-right: 5px;padding-right:8px;width: 210px;" :label="item.token" border>
                  {{ item.publisherName }} commitID: {{ item.commit }}
                  <span v-if="item.publishState === 1" style="color:#67C23A;float:right;line-height:16px;">成功</span>
                  <span v-else style="color:#F56C6C;float:right;line-height:16px;">失败</span>
                </el-radio>
                <el-button type="danger" plain @click="rollback(item)">rebuild</el-button>
              </el-row>
            </el-row>
          </el-radio-group>
        </el-col>
        <el-col :span="16" class="project-detail">
          <el-row v-for="(item, index) in publishLocalTraceList" :key="index">
            <el-row v-if="item.type === 2">
              <el-row style="margin:5px 0">git同步信息</el-row>
              <el-row style="margin:5px 0">时间: {{ item.createTime }}</el-row>
              <!-- 用数组的形式 兼容以前版本 -->
              <el-row v-if="item.state !== 0">
                <el-row>commit: {{ item['commit'] }}</el-row>
                <el-row>message: {{ item['message'] }}</el-row>
                <el-row>author: {{ item['author'] }}</el-row>
                <el-row>datetime: {{ item['timestamp'] ? parseTime(item['timestamp']) : '' }}</el-row>
                <el-row><span v-html="formatDetail(item['diff'])" /></el-row>
              </el-row>
              <el-row v-else style="margin:5px 0">
                <el-tag type="danger" effect="plain">失败</el-tag>
                <span v-html="formatDetail(item.detail)" />
              </el-row>
            </el-row>
            <el-row v-if="item.type === 3">
              <hr>
              <el-row style="margin:5px 0">获取代码后脚本信息</el-row>
              <el-row style="margin:5px 0">时间: {{ item.createTime }}</el-row>
              <el-row>脚本: <pre v-html="formatDetail(item.script)" /></el-row>
              <el-row style="margin:5px 0">
                <el-tag v-if="item.state === 0" type="danger" effect="plain">失败</el-tag>
                <el-row>[goploy ~]# <span v-html="formatDetail(item.detail)" /></el-row>
              </el-row>
            </el-row>
            <el-row v-if="item.type === 6">
              <hr>
              <el-row>
                <el-row style="margin:5px 0">remote服务器信息</el-row>
                <el-row style="margin:5px 0">服务器: {{ item.serverName }}</el-row>
                <el-row style="margin:5px 0">时间: {{ item.createTime }}</el-row>
                <el-row style="margin:5px 0">脚本: {{ item.script }}</el-row>
                <el-row style="margin:5px 0">
                  <el-tag v-if="item.state === 0" type="danger" effect="plain">失败</el-tag>
                  <span v-html="formatDetail(item.detail)" />
                </el-row>
              </el-row>
            </el-row>
          </el-row>
          <hr>
          <el-row style="margin:5px 0">目标服务器</el-row>
          <el-tabs v-model="activeRomoteTracePane">
            <el-tab-pane v-for="(item, serverName) in publishRemoteTraceList" :key="serverName" :label="serverName" :name="serverName">
              <el-row v-for="(trace, key) in item" :key="key">
                <el-row v-if="trace.type === 4">
                  <el-row style="margin:5px 0">部署前脚本</el-row>
                  <el-row style="margin:5px 0">时间: {{ trace.createTime }}</el-row>
                  <el-row>脚本: <pre v-html="formatDetail(trace.script)" /></el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">失败</el-tag>
                    <el-row>[goploy ~]# <span v-html="formatDetail(trace.detail)" /></el-row>
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 5">
                  <el-row style="margin:5px 0">rsync同步文件</el-row>
                  <el-row style="margin:5px 0">时间: {{ trace.createTime }}</el-row>
                  <el-row>命令: {{ trace.command }}</el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">失败</el-tag>
                    <span v-html="formatDetail(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else>
                  <el-row style="margin:5px 0">部署后脚本</el-row>
                  <el-row style="margin:5px 0">时间: {{ trace.createTime }}</el-row>
                  <el-row>脚本: {{ trace.script }}</el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">失败</el-tag>
                    <el-row>[goploy ~]# <span v-html="formatDetail(trace.detail)" /></el-row>
                  </el-row>
                </el-row>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="commit管理" :visible.sync="commitDialogVisible">
      <el-table
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="commitTableData"
      >
        <el-table-column type="expand">
          <template slot-scope="props">
            <span v-html="formatDetail(props.row.diff)" />
          </template>
        </el-table-column>
        <el-table-column prop="commit" label="commit" width="290" />
        <el-table-column prop="author" label="author" />
        <el-table-column label="提交时间" width="135">
          <template slot-scope="scope">
            {{ parseTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="message" label="message" width="200" show-overflow-tooltip />
        <el-table-column prop="operation" label="操作" width="80" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" @click="rollback(scope.row)">构建</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="commitDialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, getDetail, getPreview, getCommitList, publish } from '@/api/deploy'
import { getDeployOption as getDeployGroupOption } from '@/api/group'
import { parseTime } from '@/utils'

export default {
  data() {
    return {
      groupId: parseInt(localStorage.getItem('groupId')) || 0,
      groupOption: [],
      projectName: '',
      publishToken: '',
      commitDialogVisible: false,
      dialogVisible: false,
      webSocket: null,
      tableData: [],
      commitTableData: [],
      gitTraceList: [],
      publishTraceList: [],
      publishLocalTraceList: [],
      publishRemoteTraceList: {},
      activeRomoteTracePane: ''
    }
  },
  created() {
    this.getList()
    this.getDeployGroupOption()
    this.connectWebSocket()
    // // 路由跳转时结束websocket链接
    this.$router.afterEach(() => {
      this.webSocket && this.webSocket.close()
    })
  },
  methods: {
    parseTime,
    connectWebSocket() {
      if (this.webSocket && this.webSocket.readyState < 2) {
        console.log('reusing the socket connection [state = ' + this.webSocket.readyState + ']: ' + this.webSocket.url)
        return Promise.resolve(this.webSocket)
      }

      return new Promise((resolve, reject) => {
        this.webSocket = new WebSocket('ws://' + window.location.host + process.env.VUE_APP_BASE_API + '/ws/broadcast')

        this.webSocket.onopen = () => {
          console.log('socket connection is opened [state = ' + this.webSocket.readyState + ']: ' + this.webSocket.url)
          resolve(this.webSocket)
        }

        this.webSocket.onerror = (err) => {
          console.error('socket connection error : ', err)
          reject(err)
        }

        this.webSocket.onclose = (e) => {
          this.webSocket = null
          console.log('connection closed (' + e.code + ')')
        }

        this.webSocket.onmessage = (e) => {
          const data = JSON.parse(e.data)
          console.log(data)
          data.message = this.formatDetail(data.message)
          if (data.state === 3) {
            this.$notify.error({
              title: data.projectName,
              dangerouslyUseHTMLString: true,
              message: data.message,
              duration: 0
            })
          }
          const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
          if (projectIndex !== -1) {
            this.tableData[projectIndex].deployState = data.state
            this.tableData[projectIndex].publisherName = data.username
            this.tableData[projectIndex].updateTime = parseTime(new Date())
          }
        }
      })
    },

    handleGroupChange(groupId) {
      localStorage.setItem('groupId', groupId)
      this.groupId = groupId
      this.getList()
    },

    getDeployGroupOption() {
      getDeployGroupOption().then((response) => {
        this.groupOption = response.data.groupList || []
      })
    },

    getList() {
      getList(this.groupId, this.projectName).then((response) => {
        const projectList = response.data.projectList || []
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.tableData = projectList
      })
    },

    publish(data) {
      const id = data.id
      const h = this.$createElement
      let color = ''
      if (data.environment === '生产环境') {
        color = 'color: #F56C6C'
      } else if (data.environment === '测试环境') {
        color = 'color: #E6A23C'
      } else {
        color = 'color: #909399'
      }
      this.$confirm('', '提示', {
        message: h('p', null, [
          h('span', null, '此操作将部署 ' + data.name),
          h('b', { style: color }, '(' + data.environment + ')'),
          h('span', null, ', 是否继续? ')
        ]),
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.gitLog = []
        this.remoteLog = {}
        this.connectWebSocket().then(server => {
          publish(id, '').then((response) => {
            const projectIndex = this.tableData.findIndex(element => element.id === id)
            this.tableData[projectIndex].deployState = 1
          })
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消构建'
        })
      })
    },

    getDetail() {
      getDetail(this.publishToken).then((response) => {
        const publishTraceList = response.data.publishTraceList || []
        this.publishTraceList = publishTraceList.map(element => {
          element.createTime = parseTime(element.createTime)
          if (element.ext !== '') Object.assign(element, JSON.parse(element.ext))
          return element
        })

        this.publishLocalTraceList = this.publishTraceList.filter(element => element.type < 4)
        this.publishRemoteTraceList = {}
        for (const trace of this.publishTraceList) {
          if (trace.type < 4) continue
          if (!this.publishRemoteTraceList[trace.serverName]) {
            this.publishRemoteTraceList[trace.serverName] = []
          }
          this.publishRemoteTraceList[trace.serverName].push(trace)
        }
        this.activeRomoteTracePane = Object.keys(this.publishRemoteTraceList)[0]
      })
    },

    handleDetail(data) {
      this.dialogVisible = true
      getPreview(data.id).then((response) => {
        const gitTraceList = response.data.gitTraceList || []
        this.gitTraceList = gitTraceList.map(element => {
          if (element.ext !== '') Object.assign(element, JSON.parse(element.ext))
          element.commit = element['commit'] ? element['commit'].substring(0, 6) : ''
          return element
        })
        if (this.gitTraceList.length > 0) {
          this.publishToken = this.gitTraceList[0].token
          this.getDetail()
        }
      })
    },

    handleDetailChange(lastPublishToken) {
      this.publishToken = lastPublishToken
      this.getDetail()
    },

    handlePublishCommand(data) {
      const id = data.id
      getCommitList(id).then(response => {
        this.commitTableData = response.data.commitList.map(element => {
          return Object.assign(element, { projectId: id })
        })
        this.commitDialogVisible = true
      })
    },

    rollback(data) {
      this.$confirm('此操作将重新构建' + data.commit + ', 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.connectWebSocket().then(server => {
          publish(data.projectId, data.commit).then((response) => {
            const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
            this.tableData[projectIndex].deployState = 1
            this.commitDialogVisible = false
          })
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消重新构建'
        })
      })
    },

    findGroupName(groupId) {
      const projectGroup = this.groupOption.find(element => element.id === groupId)
      return projectGroup ? projectGroup['name'] : '默认'
    },

    formatDetail(detail) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import "@/styles/mixin.scss";
.project-detail {
  height:375px;
  overflow-y: auto;
  @include scrollBar();
}
.publish-btn {
  margin-right: 10px;
  >>>.el-dropdown__caret-button{
    padding-bottom: 6px;
  }
}

@media screen and (max-width: 1440px){
  .publish-record{
    >>>.el-dialog {
      width: 75%;
    }
  }
}
</style>
