# 运动接口详情

## 1.1 获取打卡记录列表

GET /api/sport/get-record

### 请求参数
格式
```json
{
  "year": "（必填）年份",
  "month": "（必填）月份",
  "date": "（可选）日期"
}
```
示例
```
/api/sport/get-record?year=2018&month=8&date=25
```

### 返回结果
```json
{
  "code": 0,
  "data": [
    { "year": 2018, "month": 8, "date": 1, "isPunch": true },
    { "year": 2018, "month": 8, "date": 2, "isPunch": false },
    { "year": 2018, "month": 8, "date": 3, "isPunch": true },
    { "year": 2018, "month": 8, "date": 4, "isPunch": true },
    { "year": 2018, "month": 8, "date": 5, "isPunch": false },
  ]
}
```

- data内返回所选范围内所有天的打卡情况


## 1.2 保存打卡记录列表

POST /api/sport/update-record

### 请求参数（body体）

格式
```json
{
  "year": "（可选）年份，默认为当年",
  "month": "（可选）月份，默认为当月",
  "date": "（可选）日期，默认为当天"
}
```

示例
```json
{
  "year": 2018,
  "month": 8,
  "date": 25
}
```
### 返回结果

```json
{
  "code": 0,
  "data": { "isPunch": true }
}
```

- isPunch表示
    - true 当日已打卡
    - false 取消当日打卡