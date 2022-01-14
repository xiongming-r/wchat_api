const express = require('express');
const fs = require('fs');//引入文件读取模块
const cors =require('cors')
// const Path =require('path');
const NODE_ENV = process.env.NODE_ENV
const { exec } = require('child_process')
// const { set } = require('core-js/core/dict');

// const bodyParser =require('body-parser')
// const homeApi=require('./api/homeApi')
// const { Nuxt, Builder } = require('nuxt');
const app = express();
app.use(cors())
app.use(express.urlencoded({ limit:'50mb',extended: false }))
app.use(express.json({limit:'50mb'}))
app.use(function(err,req,res,next){
  console.error('error:'+err.stack);

})
process.on('uncaughtException', function (err) {
  console.log('uncaughtException :'+err.stack)
  // sendErrCourier(err.stack)
})
process.on('unhandledRejection', (reason, p) => {
  console.log('Unhandled Rejection at: Promise', p, 'reason:', reason)
  // sendErrCourier(reason)
})
app.get('/estimate',function(req,res){
  console.log(req.query);
    exec('D:/code/tools/go/wayz1103.exe --action estimate', (error, stdout, stderr) => {
      console.log(stdout, stderr)
      res.send(stdout)
    })
})
app.get('/create',function(req,res){
  console.log(req.query);
    exec('D:/code/tools/go/wayz1103.exe --action create', (error, stdout, stderr) => {
      console.log(stdout, stderr)
      res.send(stdout)
    })
})
app.get('/auth',function(req,res){
  console.log(req.query);
  let url = `https://lbi-api.newayz.com/openapi/v1/precisionMarketing/threeParty/tencent/getAuthorizeUrl?scope=${req.query.scope}&advertiserId=${req.query.advertiserId}&account_type=${req.query.account_type}`
  let shell = `D:/code/tools/go/wayz1103.exe --action auth --auth-URL ${url}`
  console.log(shell);
    // exec(shell, (error, stdout, stderr) => {
    //   console.log(stdout, stderr)
    //   res.send(stdout)
    // })
})
app.post('/upload',function(req,res){
  console.log(req.body);
  let url = `https://lbi-api.newayz.com/openapi/v1/precisionMarketing/threeParty/tencent/uploadToMediaForClient`
  let shell = `D:/code/tools/go/wayz1103.exe --action upload --upload-URL ${url} --crowd ${req.body.wayzCrowdId} --advertiser-id ${req.body.advertiserId}`
  console.log(shell);
    // exec(shell, (error, stdout, stderr) => {
    //   console.log(stdout, stderr)
    //   res.send(stdout)
    // })
})
app.post('/query',function(req,res){
  console.log(req.body);
  let url = `https://lbi-api.newayz.com/openapi/v1/precisionMarketing/threeParty/tencent/query?advertiserId=${req.body.advertiserId}&wayzCrowdId=${req.body.wayzCrowdId}`
  let shell = `D:/code/tools/go/wayz1103.exe --action query --query-URL ${url}`
  console.log(shell);
    // exec(shell, (error, stdout, stderr) => {
    //   console.log(stdout, stderr)
    //   res.send(stdout)
    // })
})



// app.use('/home',homeApi)

app.listen(3001)
console.log('success listen at port:3001');



console.log('服务器开启中...');


