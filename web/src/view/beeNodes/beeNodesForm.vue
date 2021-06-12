<template>
<div>
    <el-form :model="formData" label-position="right" label-width="80px">
             <el-form-item label="uuid字段:">
                <el-input v-model="formData.uuid" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="name字段:">
                <el-input v-model="formData.name" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="ip字段:">
                <el-input v-model="formData.ip" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="debugPort字段:"><el-input v-model.number="formData.debugPort" clearable placeholder="请输入"></el-input>
          </el-form-item>
           
             <el-form-item label="walletAddress字段:">
                <el-input v-model="formData.walletAddress" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           
             <el-form-item label="uncashedCount字段:"><el-input v-model.number="formData.uncashedCount" clearable placeholder="请输入"></el-input>
          </el-form-item>
           
             <el-form-item label="peerCount字段:"><el-input v-model.number="formData.peerCount" clearable placeholder="请输入"></el-input>
          </el-form-item>
           
             <el-form-item label="ethBalance字段:">
                  <el-input-number v-model="formData.ethBalance" :precision="2" clearable></el-input-number>
           </el-form-item>
           
             <el-form-item label="bzzBalance字段:">
                  <el-input-number v-model="formData.bzzBalance" :precision="2" clearable></el-input-number>
           </el-form-item>
           
             <el-form-item label="用户id:">
                <el-input v-model="formData.user_id" clearable placeholder="请输入" ></el-input>
          </el-form-item>
           <el-form-item>
           <el-button @click="save" type="primary">保存</el-button>
           <el-button @click="back" type="primary">返回</el-button>
           </el-form-item>
    </el-form>
</div>
</template>

<script>
import {
    createBeeNodes,
    updateBeeNodes,
    findBeeNodes
} from "@/api/beeNodes";  //  此处请自行替换地址
import infoList from "@/mixins/infoList";
export default {
  name: "BeeNodes",
  mixins: [infoList],
  data() {
    return {
      type: "",formData: {
            uuid:"",
            name:"",
            ip:"",
            debugPort:0,
            walletAddress:"",
            uncashedCount:0,
            peerCount:0,
            ethBalance:0,
            bzzBalance:0,
            user_id:"",
            
      }
    };
  },
  methods: {
    async save() {
      let res;
      switch (this.type) {
        case "create":
          res = await createBeeNodes(this.formData);
          break;
        case "update":
          res = await updateBeeNodes(this.formData);
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
      }
    },
    back(){
        this.$router.go(-1)
    }
  },
  async created() {
   // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if(this.$route.query.id){
    const res = await findBeeNodes({ ID: this.$route.query.id })
    if(res.code == 0){
       this.formData = res.data.rebeeNodes
       this.type = "update"
     }
    }else{
     this.type = "create"
   }
  
}
};
</script>

<style>
</style>