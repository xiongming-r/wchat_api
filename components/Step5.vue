<template>
  <div class="main">
    <a-card title="创建" class="step1">
      <a-row :gutter="16">
        <a-col :span="4" >
          <a-form-model  ref="ruleForm" :model="form">
            <a-form-model-item ref="wayzCrowdId" label="name" prop="name">
              <a-input v-model="form.name"></a-input>
            </a-form-model-item>
            <a-form-model-item ref="wayzCrowdId" label="desc" prop="desc">
              <a-input v-model="form.desc"></a-input>
            </a-form-model-item>
            <a-form-model-item ref="tagCodes" label="tagCodes" prop="tagCodes">
              <a-input v-model="form.tagCodes"></a-input>
            </a-form-model-item>
            <a-form-model-item :wrapper-col="{ span: 14, offset: 4 }">
              <a-button type="primary" @click="onSubmit">
                创建人群包
              </a-button>
            </a-form-model-item>
          </a-form-model>
        </a-col>
        <a-col :span="20">
          <codemirror 
            ref="safetyCmEditor"
            :value="safetyCode"
            :options="cmOptions"
            @ready="onCmReady"
          />
        </a-col>
      </a-row>
    </a-card>
  </div>
</template>

<script>
import api from '@/api/baseApi.js'
export default {
  data(){
    return{
      safetyCode:'',
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      form:{
        name:'',
        tagCodes:'',
        desc:'',
      },
      cmOptions: {
        tabSize: 4,
        // mode: 'text/javascript',
        theme: 'base16-light',
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
      }
    }
  },
  methods:{
    onChange(e){
      this.form.scope=e.target.value
    },
    async onSubmit(){
      try {
        let params={
          name:this.form.name,
          desc:this.form.desc,
          tagCodes:[this.form.tagCodes],
        }
        await this.$axios.post('http://localhost:3001/create',params).then(res=>{
          console.log(res.data);
          // this.safetyCode=JSON.stringify(res.data,null,'\t')
          this.safetyCode=res.data
        })
      } catch (error) {
        console.log(error);
        this.safetyCode=JSON.stringify(error,null,'\t')
      }
    },
    resetForm(){
      this.$refs.ruleForm.resetFields();
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