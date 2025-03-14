
执行命令
```
chatgpt_api company_list_aa.csv 
或
nohup ./chatgpt_api company_list_aa.csv >> company_list_aa.log 2>&1 &
```

config说明

```
config/gdbc.yaml文件：

# 关键词的前缀、后缀
batch_ask_prefix: "请帮我判断一下,"
batch_ask_suffix: ", 基于国民经济行业分类（GB/T 4754-2017)标准下钻到四级分类, 不加说明仅返回|分割的一、二、三、四级分类名称"

```

目录文件如下：

```
config
chatgpt_api.exe
sample_company_1.csv
sample_company_1_output_from_ai.csv
```

其中company_list_aa.csv文件内容：
```
行业一级分类,公司名
G,东部机场集团有限公司
G,中国东方航空股份有限公司
```

输入输出:

```
Question:
请帮我判断一下,行业一级分类:C, 企业名称:希捷国际科技(无锡)有限公司, 基于国民经济行业分类（GB/T 4754-2017)标准下钻到四级分类, 不加说明仅返回|分割的一、二、三、四级分类名称
Answer:
 C|制造业|计算机、通信和其他电子设备制造业|存储设备制造
```

程序日志示例：

```
2025/03/14 12:22:22   --> req json: https://api.deepseek.com/chat/completions {"model":"deepseek-chat","messages":[{"role":"user","content":"请帮我判断一下,行业一级分类:C, 企业名称:希捷国际科技(无锡)有限公司, 基于国民经济行业分类（GB/T 4754-2017)标准下钻到四级分类, 不加说明仅返回|分割的一、二、三、四级分类名称"}],"temperature":0}
2025/03/14 12:22:27   <-- 回复: C|制造业|计算机、通信和其他电子设备制造业|存储设备制造

```

最终生成文件company_list_aa_output_from_ai.csv如下：

```
行业一级分类:G， 公司名:东部机场集团有限公司,G,交通运输、仓储和邮政业,航空运输业,机场,民用机场
行业一级分类:G， 公司名:中国东方航空股份有限公司,G,交通运输、仓储和邮政业,航空运输业,航空客货运输
```
