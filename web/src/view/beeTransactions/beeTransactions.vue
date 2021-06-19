<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
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
    <el-table-column label="取票日期" width="160">
         <template slot-scope="scope">{{scope.row.CreatedAt|formatDate}}</template>
    </el-table-column>
    <el-table-column label="ip" prop="ip" width="120"></el-table-column>
    <el-table-column label="peer" prop="peer" width="250"></el-table-column>
    <el-table-column label="钱包地址" min-width="150">
        <template slot-scope="scope">
          <a :href="'https://goerli.etherscan.io/address/'+scope.row.walletAddress"
            target="_blank" style="color:#2665e4">{{scope.row.walletAddress}}
          </a>
        </template>
    </el-table-column>
    <el-table-column label="txid" min-width="150">
        <template slot-scope="scope">
          <a :href="'https://goerli.etherscan.io/tx/'+scope.row.txid"
            target="_blank" style="color:#2665e4">{{scope.row.txid}}
          </a>
        </template>
    </el-table-column>
    <el-table-column label="gas price" prop="gasPrice" width="80"></el-table-column>
    <el-table-column label="领取票值" prop="cashedAmount" width="80"></el-table-column>
    <el-table-column label="nonce" prop="nonce" width="100"></el-table-column>
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
  </div>
</template>

<script>
import {
    getBeeTransactionList,
} from "@/api/beeTransactions";  //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/date";
import infoList from "@/mixins/infoList";
export default {
  name: "BeeTransactions",
  mixins: [infoList],
  data() {
    return {
      listApi: getBeeTransactionList,
      multipleSelection: [],
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
  },
  async created() {
    await this.getTableData();
  }
};
</script>

<style>
</style>
