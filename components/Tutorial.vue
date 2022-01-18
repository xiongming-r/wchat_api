<template>
  <div class="main">
    <Step4 />
    <Step5 />
    <a-card title="授权" class="step1">
      <a-row :gutter="16">
        <a-col :span="12" >
          <a-form-model  ref="ruleForm" :model="form">
            <a-form-model-item ref="scope" label="scope" prop="scope">
              <a-radio-group @change="onChange">
                <a-radio-button v-for="(item,index) in list" :key="index" :value="item.value">
                  {{item.name}}
                </a-radio-button>
              </a-radio-group>
              <a-input v-model="form.scope"></a-input>
            </a-form-model-item>
            <a-form-model-item ref="advertiserId" prop="advertiserId" label="advertiserId">
              <a-input v-model="form.advertiserId"></a-input>
            </a-form-model-item>
            <a-form-model-item ref="account_type" prop="account_type" label="account_type">
              <a-select v-model="form.account_type">
                <a-select-option v-for="(item,index) in listForaccounttype" :key="index" :value="item.value">{{item.value}}</a-select-option>
              </a-select>
            </a-form-model-item>
            <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
              <a-button type="primary" @click="onSubmit">
                提交
              </a-button>
              <a-button style="margin-left: 10px;" @click="resetForm">
                重置
              </a-button>
              <!-- <nuxt-link v-if="authSuccs" :to="http">点击复制网址</nuxt-link> -->
              <a-button v-if="authSuccs" style="margin-left: 10px;" v-clipboard:copy="http" v-clipboard:success="onCopy" v-clipboard:error="onError">
                点击复制网址
              </a-button>
            </a-form-model-item>
          </a-form-model>
        </a-col>
        <a-col :span="12">
          <codemirror 
            ref="safetyCmEditor"
            :value="safetyCode"
            :options="cmOptions"
            @ready="onCmReady"
          />
        </a-col>
      </a-row>
    </a-card>
    <Step2 />
    <Step3 />
  </div>
</template>

<script>
import api from '@/api/baseApi.js'
import Step2 from '@/components/Step2.vue'
import Step3 from '@/components/Step3.vue'
import Step4 from '@/components/Step3.vue'
import Step5 from '@/components/Step3.vue'
export default {
  components:{
    Step2,
    Step3
  },
  data(){
    return{
      list:[
        {
          name:'广告投放',
          value:'ads_management'
        },
        {
          name:'数据洞察',
          value:'ads_insights'
        },
        {
          name:'账号服务',
          value:'account_management'
        },
        {
          name:'人群管理',
          value:'audience_management'
        },
        {
          name:'用户行为数据接入',
          value:'user_actions'
        },
      ],
      listForaccounttype:[
        {
          name:'ACCOUNT_TYPE_WECHAT',
          value:'ACCOUNT_TYPE_WECHAT'
        },
        {
          name:'ACCOUNT_TYPE_QQ',
          value:'ACCOUNT_TYPE_QQ'
        }
      ],
      safetyCode:'',
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      form:{
        scope:'',
        advertiserId:'',
        account_type:''
      },
      cmOptions: {
        tabSize: 4,
        // mode: 'text/javascript',
        theme: 'elegant',
        lineNumbers: true,
        line: true,
        styleActiveLine: true,
        highlightDifferences: true,
        mode: { // 模式, 可查看 codemirror/mode 中的所有模式
          name: 'javascript',
          json: true
        },
        hintOptions: {
          // 当匹配只有一项的时候是否自动补全
          completeSingle: false
        }, 
        // 快捷键
        keyMap:'sublime',
        extraKeys:{'Ctrl':'autocomplete'}
        // more CodeMirror options...
      },
      authSuccs:false,
      http:''
    }
  },
  mounted(){
    // this.$axios.get('http://localhost:3001/getPath').then(res=>{
    //   console.log(res,'1111');
    // })
  },
  methods:{
    onChange(e){
      this.form.scope=e.target.value
    },
    async onSubmit(){
      try {
        // let url = `https://lbi-api.newayz.com/openapi/v1/precisionMarketing/threeParty/tencent/getAuthorizeUrl?scope=${this.form.scope}&advertiserId=${this.form.advertiserId}&account_type=${this.form.account_type}` 
        // console.log(url);
        await this.$axios.get('http://localhost:3001/auth'+`?scope=${this.form.scope}&advertiserId=${this.form.advertiserId}&account_type=${this.form.account_type}`).then(res=>{
          console.log(res);
          console.log();
          let http = this.httpString(res.data)
          let index =res.data.lastIndexOf('resp=')
          let obj = res.data.substring(index+5,res.data.length)
          this.safetyCode=JSON.stringify(obj,null,'\t')
          if(http.length>0){
            this.authSuccs=true
            this.http=http[0]
          }
        })
      } catch (error) {
        console.log(error);
        this.safetyCode=JSON.stringify(error,null,'\t')
      }
    },
    goHttp(){

    },
    onCopy(e){
      this.$message.success('已复制')
    },
    onError(e){
      this.$message.error('复制失败')
    },
    resetForm(){
      this.$refs.ruleForm.resetFields();
    },
    httpString(s) {
     var reg = /(https?|http|ftp|file):\/\/[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]/g;
      s = s.match(reg);
     return (s)
    },
    onCmReady(cm) {
      console.log('the editor is readied!', cm)
      // setTimeout(() => {
      //   this.$refs.cmEditor.refresh()
      //   console.log(this.$refs.cmEditor,'222');
      // }, 100);
      cm.on('inputRead',(cm,obj)=>{
        if(obj.text && obj.text.length>0){
          const c = obj.text[0].charAt(obj.text[0].length-1)
          if((c>='a' && c <= 'z')|| (c>='A' && c<= 'Z')){
            cm.showHint({completeSingle:false})
          }
        }
      })
    },
    onCmFocus(cm) {
      console.log('the editor is focused!', cm)
      cm.on('inputRead',(cm,obj)=>{
        if(obj.text && obj.text.length>0){
          const c = obj.text[0].charAt(obj.text[0].length-1)
          if((c>='a' && c <= 'z')|| (c>='A' && c<= 'Z')){
            cm.showHint({completeSingle:false})
          }
        }
      })
    },
  }
}
</script>

<style lang="less" scoped>

</style>