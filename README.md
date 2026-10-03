# 学业审核工具（course-credit-audit）

仅使用 Go 标准库实现的毕业学业审核工具。读取 JSON 任务文件，
按课程目录与毕业规则评估每个学生的修读记录，输出 JSON 报告，
说明学分缺口与先修违规。不连接任何教务系统。

## 运行

```sh
go run ./cmd/degree <task.json> <report.json>
```

任一输入非法时以非零状态退出，且不生成报告文件。

## 任务格式

```json
{
  "courses": [
    {"id": "CS101", "credits": 3, "category": "cs", "prerequisites": []}
  ],
  "requirements": {
    "required": ["CS101"],
    "categoryMin": {"cs": 6},
    "totalMin": 10
  },
  "students": [
    {"id": "S001", "records": [
      {"course": "CS101", "term": 1, "passed": true}
    ]}
  ]
}
```

- 课程体系：每门课程编号唯一，`credits` 为正整数，`category` 非空，
  `prerequisites` 为直接先修课程编号。先修必须存在，拒绝自依赖与环路。
- 毕业规则：`required` 为必修课程集合，`categoryMin` 为各类别最低学分，
  `totalMin` 为总学分门槛，门槛均为非负整数，未知课程或类别将被拒绝。
- 修读记录：`term` 为正整数学期序号；同一学生同课程同学期不能重复。

## 审核语义

- 每门课只在首次及格时取得一次学分；之后不及格不扣回，重复及格不重复累计。
- 每次修读要求全部直接先修曾在严格更早学期及格；同学期通过不满足先修，
  未来成绩不影响过去的判断。违规逐条列出课程、学期与缺失先修。
- 先修违规不撤销已取得的学分，但存在违规即不能判定毕业。
- 一门课同时计入所属类别学分与总学分；必修全部完成、各类别与总学分
  门槛均满足且无先修违规时 `graduated` 为 `true`。
- 空成绩记录正常报告全部缺口。

## 报告格式

按学生编号稳定排序输出；`missingRequired` 按课程编号排序，
`categoryDeficits` 仅包含非负缺口（缺口为 0 的类别不列出）：

```json
{
  "students": [
    {
      "id": "S001",
      "totalCredits": 10,
      "categoryCredits": {"cs": 6, "math": 4},
      "missingRequired": [],
      "categoryDeficits": {},
      "totalDeficit": 0,
      "prerequisiteViolations": [
        {"course": "CS201", "term": 1, "missing": ["CS101"]}
      ],
      "graduated": false
    }
  ]
}
```

## 测试

```sh
go test ./...
```

覆盖重修（首次及格取得学分、后续不及格不扣回、重复及格不重复累计）、
同学期先修违规、类别达标而总量不足、空成绩缺口、环路/自依赖/未知引用
等非法输入，以及输出的稳定排序。

示例任务见 `testdata/task.json`。
