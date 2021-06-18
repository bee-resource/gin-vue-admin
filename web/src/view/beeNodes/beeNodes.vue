<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item>
          <el-button @click="openImportBeeNodesDialog" type="primary">批量新增节点</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="openCashoutDialog" type="success">批量取票</el-button>
        </el-form-item>

        <el-form-item>
          <el-button @click="batchRefresh" type="warning">批量查看节点状态</el-button>
        </el-form-item>

        <el-form-item>
          <el-popover placement="top" v-model="deleteVisible" width="160">
            <p>确定要删除吗？</p>
              <div style="text-align: right; margin: 0">
                <el-button @click="deleteVisible = false" size="mini" type="text">取消</el-button>
                <el-button @click="onDelete" size="mini" type="primary">确定</el-button>
              </div>
            <el-button icon="el-icon-delete" size="mini" slot="reference" type="danger">批量删除</el-button>
          </el-popover>
        </el-form-item>
      </el-form>
    </div>
    <el-table
      :data="tableData"
      @selection-change="handleSelectionChange"
      border
      ref="multipleTable"
      stripe
      style="width: 100%"
      tooltip-effect="dark"
    >
    <el-table-column type="selection" width="40"></el-table-column>
    <el-table-column label="上次更新日期" width="160">
         <template slot-scope="scope">{{scope.row.UpdatedAt|formatDate}}</template>
    </el-table-column>

    <el-table-column label="名称" prop="name" width="60"></el-table-column>

    <el-table-column label="版本" prop="version" width="85">
    </el-table-column>

    <el-table-column label="ip" prop="ip" width="120"></el-table-column>

    <el-table-column label="端口" prop="debugPort" width="60"></el-table-column>

    <el-table-column label="钱包地址" min-width="100">
        <template slot-scope="scope">
          <a :href="'https://goerli.etherscan.io/address/'+scope.row.walletAddress"
            target="_blank">{{scope.row.walletAddress}}
          </a>
        </template>
    </el-table-column>

    <el-table-column label="未领票数" prop="uncashedCount" width="80"></el-table-column>
    <el-table-column label="未领票值" prop="uncashedAmount" width="80"></el-table-column>
    <el-table-column label="连接数" prop="peerCount" width="70"></el-table-column>
    <el-table-column label="eth余额" prop="ethBalance" width="100"></el-table-column>
    <el-table-column label="bzz余额" prop="bzzBalance" width="100"></el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="scope">
          <el-button class="table-button" @click="updateBeeNodesStatus(scope.row)" size="small" type="warning">查</el-button>
          <el-button class="table-button" @click="cashoutBeeNodes(scope.row)" size="small" type="success">收</el-button>
          <el-button class="table-button" @click="updateBeeNodes(scope.row)" size="small" type="primary">改</el-button>
          <el-button type="danger" size="mini" @click="deleteRow(scope.row)">删</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{float:'right',padding:'20px'}"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>

    <el-dialog :before-close="closeBatchImportDialog" :visible.sync="batchImportDialogFormVisible">
      <el-form :model="batchImportFormData" label-position="right" label-width="80px">
        <el-form-item label="批量导入">
          <el-input v-model="batchImportFormData.nodes" type="textarea"
            placeholder="可导入多个节点，每行为一个节点，格式为ip:port，或者只输入ip，port默认为1635批量导入"
            :autosize="{minRows: 20, maxRows: 500}" :style="{width: '100%'}"></el-input>
        </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeBatchImportDialog">取消</el-button>
        <el-button type="primary" @click="enterBatchImportDialog">确定</el-button>
      </div>
    </el-dialog>

    <el-dialog :before-close="closeDialog" :visible.sync="dialogFormVisible" title="弹窗操作">
      <el-form :model="formData" label-position="right" label-width="80px">
         <el-form-item label="名称">
            <el-input v-model="formData.name" clearable placeholder="请输入" ></el-input>
          </el-form-item>
          <el-form-item label="ip">
              <el-input v-model="formData.ip" clearable placeholder="请输入" ></el-input>
          </el-form-item>
            <el-form-item label="端口"><el-input v-model.number="formData.debugPort" clearable placeholder="请输入"></el-input>
          </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button @click="enterDialog" type="primary">确 定</el-button>
      </div>
    </el-dialog>

    <el-dialog :before-close="closeCashOutDialog" :visible.sync="cashoutDialogFormVisible" title="弹窗操作">
      <el-form :model="cashoutFormData" label-position="right" label-width="80px">
          <el-form-item label="gasPrice(单位G)"><el-input v-model.number="cashoutFormData.gasPrice" clearable placeholder="单位为G，默认800"></el-input>
          </el-form-item>
          <el-form-item label="count">
              <el-input v-model="cashoutFormData.count" clearable placeholder="请输入取票次数" ></el-input>
          </el-form-item>
            <el-form-item label="nonce"><el-input v-model.number="cashoutFormData.nonce" clearable placeholder="不理解请保持-1即可"></el-input>
          </el-form-item>
      </el-form>
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeCashOutDialog">取 消</el-button>
        <el-button @click="enterCashoutDialog" type="primary">确 定</el-button>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import {
    createBeeNodes,
    deleteBeeNodes,
    deleteBeeNodesByIds,
    updateBeeNodes,
    updateBeeNodesStatus,
    findBeeNodes,
    getBeeNodesList,
    importBeeNodes,
    updateBeeNodeStatusInBatch,
    cashoutBeeNodes
} from "@/api/beeNodes";  //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
export default {
  name: "BeeNodes",
  mixins: [infoList],
  data() {
    return {
      listApi: getBeeNodesList,
      dialogFormVisible: false,
      batchImportDialogFormVisible: false,
      cashoutDialogFormVisible: false,
      type: "",
      deleteVisible: false,
      multipleSelection: [],
      formData: {
            name:"",
            ip:"",
            debugPort:1635,
      },
      batchImportFormData : {
        nodes: ""
      },
      cashoutFormData: {
        gasPrice: 800,
        count: 1,
        nonce: -1,
      }
    };
  },
  filters: {
    formatDate: function(time) {
      if (time != null && time != "") {
        var date = new Date(time);
        return formatTimeToStr(date, "yyyy-MM-dd hh:mm:ss");
      } else {
        return "";
      }
    },
    formatBoolean: function(bool) {
      if (bool != null) {
        return bool ? "是" :"否";
      } else {
        return "";
      }
    }
  },
  methods: {
      //条件搜索前端看此方法
      onSubmit() {
        this.page = 1
        this.pageSize = 10
        this.getTableData()
      },
      handleSelectionChange(val) {
        this.multipleSelection = val
      },
      deleteRow(row){
        this.$confirm('确定要删除吗?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
           this.deleteBeeNodes(row);
        });
      },
      async onDelete() {
        const ids = []
        if(this.multipleSelection.length == 0){
          this.$message({
            type: 'warning',
            message: '请选择要删除的数据'
          })
          return
        }
        this.multipleSelection &&
          this.multipleSelection.map(item => {
            ids.push(item.ID)
          })
        const res = await deleteBeeNodesByIds({ ids })
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功'
          })
          if (this.tableData.length == ids.length && this.page > 1) {
              this.page--;
          }
          this.deleteVisible = false
          this.getTableData()
        }
      },
    async batchRefresh() {
      const ids = []
        if(this.multipleSelection.length == 0){
          this.$message({
            type: 'warning',
            message: '请选择要查看状态的节点'
          })
          return
        }
      this.multipleSelection &&
          this.multipleSelection.map(item => {
            ids.push(item.ID)
          })
      const res = await updateBeeNodeStatusInBatch({ ids })
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '批量查看状态成功'
          })
          this.deleteVisible = false
          this.getTableData()
        }
    },
    async updateBeeNodes(row) {
      const res = await findBeeNodes({ ID: row.ID });
      this.type = "update";
      if (res.code == 0) {
        this.formData = res.data.rebeeNodes;
        this.dialogFormVisible = true;
      }
    },
    async cashoutBeeNodes(row) {
      this.cashoutFormData.current_row = row;
      this.cashoutDialogFormVisible = true;
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
          uuid:"",
          name:"",
          version:"",
          ip:"",
          debugPort:0,
          walletAddress:"",
          uncashedCount:0,
          uncashedAmount:0,
          peerCount:0,
          ethBalance:0,
          bzzBalance:0,
          user_id:"",
      };
    },
    closeBatchImportDialog() {
      this.batchImportDialogFormVisible = false;
      this.batchImportFormData = {
        nodes: ""
      };
    },

    closeCashOutDialog() {
      this.cashoutDialogFormVisible = false;
      this.cashoutFormData = {
        ids: [],
        gasPrice: 800,
        count: 1,
        nonce: -1,
      };
    },
    async deleteBeeNodes(row) {
      const res = await deleteBeeNodes({ ID: row.ID, userId: row.userId });
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "删除成功"
        });
        if (this.tableData.length == 1 && this.page > 1 ) {
            this.page--;
        }
        this.getTableData();
      }
    },
    async updateBeeNodesStatus(row) {
      const res = await updateBeeNodesStatus({ ID: row.ID, userId: row.userId });
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "更新状态成功"
        });
        this.getTableData();
      }
    },
    async enterDialog() {
      let res;
      switch (this.type) {
        case "create":
          res = await createBeeNodes(this.formData);
          break;
        case "update":
          res = await updateBeeNodes(this.formData);
          break;
        case "refresh":
          res = await updateBeeNodesStatus(this.formData);;
          break;
        default:
          res = await createBeeNodes(this.formData);
          break;
      }
      if (res.code == 0) {
        this.$message({
          type:"success",
          message:"创建/更改成功"
        })
        this.closeDialog();
        this.getTableData();
      }
    },
    async enterBatchImportDialog() {
      const ipPortList = []
      let nodes = this.batchImportFormData.nodes.split("\n");
      for (let index = 0; index < nodes.length; index++) {
        const ipPort = nodes[index].split(":");
        let port = 1635;
        if(ipPort.length == 2) {
          port = ipPort[1];
        }
        ipPortList.push({Ip: ipPort[0], Port: parseInt(port)})
      }
      let res = await importBeeNodes({ipPortList: ipPortList});
      if (res.code == 0) {
        this.$message({
          type:"success",
          message:"创建/更改成功"
        })
        this.closeBatchImportDialog();
        this.getTableData();
      }
    },
    async enterCashoutDialog() {
      const cashoutList = [];
      let gasPrice = (this.cashoutFormData.gasPrice * 10 ** 9).toString();
      let current_row_count = parseInt(this.cashoutFormData.current_row.count);
      if (current_row_count > 0) {
        cashoutList.push({ Id: this.cashoutFormData.current_row.ID, Count: max(this.cashoutFormData.count, current_row_count), GasPrice: gasPrice })
      }
      this.multipleSelection &&
        this.multipleSelection.map(item => {
          let row_count = parseInt(item.count)
          if(row_count > 0) {
            cashoutList.push( { Id: item.ID, Count: max(this.cashoutFormData.count, row_count), GasPrice: gasPrice } )
          }
        })

      if(cashoutList.length == 0) {
        this.$message({
          type:"warning",
          message:"请先选中节点或者当前选中节点无票"
        })
        return;
      }
      const res = await cashoutBeeNodes({ cashoutList: cashoutList })
      if (res.code == 0) {
        this.$message({
          type: 'success',
          message: '取票成功'
        })
        this.cashoutDialogFormVisible = false;
        this.getTableData()
      }
    },
    async cashoutDialog() {
      let count = this.cashoutFormData.count;
      let gasPrice = this.cashoutFormData.gasPrice;
      let nonce = this.cashoutFormData.nonce;
      let res = await cashoutBeeNodes({count: count, gasPrice: gasPrice, nonce: nonce});
    },
    async openImportBeeNodesDialog() {
      this.batchImportDialogFormVisible = true;
    },
    async openCashoutDialog() {
      this.cashoutDialogFormVisible = true;
    }
  },
  async created() {
    await this.getTableData();
  }
};
</script>

<style>
a:link {
  color: #2665e4;
}
</style>
