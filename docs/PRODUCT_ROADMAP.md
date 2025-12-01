# Paper AI 产品规划文档

**文档版本**: v1.0
**创建日期**: 2025-11-28
**产品定位**: 专注学术论文润色的垂直领域AI工具

---

## 📋 目录

1. [产品愿景](#1-产品愿景)
2. [当前功能分析](#2-当前功能分析)
3. [核心功能优化设计](#3-核心功能优化设计)
4. [新功能规划](#4-新功能规划)
5. [技术架构优化](#5-技术架构优化)
6. [商业模式设计](#6-商业模式设计)
7. [产品建议](#7-产品建议)
8. [实施路线图](#8-实施路线图)
9. [成功指标](#9-成功指标)

---

## 1. 产品愿景

### 1.1 核心理念

**做一件事，做到极致**

放弃"大而全"的平台模式，专注于学术论文润色这一细分领域，通过极致的产品体验建立竞争壁垒。

### 1.2 目标用户

- **主要用户**: 硕士、博士研究生
- **次要用户**: 科研工作者、高校教师
- **潜在用户**: 科研机构、实验室团队

### 1.3 用户痛点

1. **通用AI工具不专业**: ChatGPT/Claude 虽能润色，但不理解学术规范
2. **现有工具不够学术**: Grammarly 面向通用英文，缺乏学科专业性
3. **效率低下**: 逐句人工润色耗时巨大，一篇论文需要数天
4. **缺乏对比**: 无法直观看到修改前后差异和改进点
5. **学习价值低**: 工具只是修改，不解释为什么，无法提升写作能力

### 1.4 价值主张

> "让你的论文达到顶级期刊的写作水准，同时提升你的学术英文写作能力"

**三大核心价值**:
- ✅ **专业性**: 理解学科术语和学术规范
- ✅ **高效性**: 批量处理，几分钟完成全文润色
- ✅ **学习性**: 不只是改，更要让用户知道为什么这样改

---

## 2. 当前功能分析

### 2.1 已实现功能

#### 基础润色能力
```go
// 当前支持的参数
type PolishRequest struct {
    Content  string  // 内容（最大10000字符）
    Provider string  // AI提供商
    Style    string  // 风格：academic/formal/concise
    Language string  // 语言：en/zh
}
```

**优势**:
- ✅ 支持多AI提供商（可切换不同引擎）
- ✅ 基础的风格控制
- ✅ 完整的数据持久化（PostgreSQL）
- ✅ 用户认证和权限控制

**不足**:
- ❌ 只有输入框单一交互方式
- ❌ 无对比展示功能
- ❌ 无批量处理
- ❌ 风格选项过于简单
- ❌ 无修改说明和学习价值
- ❌ 缺少前端界面

### 2.2 数据库支持

已有完整的记录系统：
```go
type PolishRecord struct {
    TraceID         string  // 追踪ID
    UserID          int64   // 用户ID
    OriginalContent string  // 原文
    PolishedContent string  // 润色后
    Style, Language string  // 风格、语言
    Provider, Model string  // 提供商、模型
    ProcessTimeMs   int     // 处理时间
    Status          string  // 状态
}
```

**可以利用的数据**:
- 用户历史记录分析
- 质量统计和改进
- A/B测试不同Prompt效果

---

## 3. 核心功能优化设计

### 3.1 对比展示系统（Priority: P0）

#### 设计目标
让用户清晰看到每一处修改，理解修改价值

#### 界面设计

```
┌─────────────────────────────────────────────────────────────┐
│                      论文润色对比                            │
├──────────────────────────┬──────────────────────────────────┤
│       原文 (385词)        │      润色后 (392词)              │
├──────────────────────────┼──────────────────────────────────┤
│ In this paper, we        │ In this paper, we propose a      │
│ propose a new method     │ novel methodology for [新词:    │
│ to solve the problem.    │ methodology替代method]           │
│                          │                                  │
│ The method is based on   │ This approach is predicated on   │
│ [中式英语标记]            │ [更学术的表达]                   │
│                          │                                  │
│ [点击查看修改理由]        │ [接受] [拒绝] [查看替代方案]     │
└──────────────────────────┴──────────────────────────────────┘

统计信息:
📊 学术性提升: +23%
✏️ 修改位置: 15处
📝 词汇优化: 8处 | 语法修正: 4处 | 结构调整: 3处
```

#### 技术实现

```go
// 对比信息结构
type ComparisonResult struct {
    Original        string              `json:"original"`
    Polished        string              `json:"polished"`
    Differences     []Difference        `json:"differences"`
    Statistics      ComparisonStats     `json:"statistics"`
    Improvements    []Improvement       `json:"improvements"`
}

type Difference struct {
    Type         string   `json:"type"`         // word/phrase/sentence
    OriginalText string   `json:"original_text"`
    PolishedText string   `json:"polished_text"`
    Position     Position `json:"position"`
    Category     string   `json:"category"`     // vocabulary/grammar/style
    Reason       string   `json:"reason"`       // 修改理由
    Confidence   float64  `json:"confidence"`   // 置信度
}

type Improvement struct {
    Type        string  `json:"type"`         // academic_level/readability/clarity
    BeforeScore float64 `json:"before_score"`
    AfterScore  float64 `json:"after_score"`
    Description string  `json:"description"`
}
```

### 3.2 多版本润色（Priority: P0）

#### 设计理念
不同用户对润色"强度"有不同偏好，提供3个版本让用户选择

#### 版本定义

**保守版（Conservative）**
- 只修正明确的语法错误
- 保持原文表达方式
- 适合：已经比较好的论文，只需微调

**平衡版（Balanced）** - 默认推荐
- 适度优化词汇和句式
- 提升学术性的同时保留作者风格
- 适合：大多数情况

**激进版（Aggressive）**
- 大幅改写，追求最高学术水准
- 可能改变句式结构
- 适合：初稿质量较差，需要大幅提升

#### Prompt策略

```python
# 保守版 Prompt
"""
You are a careful academic editor. Only fix clear grammar errors
and improve obvious issues. Preserve the author's writing style
and sentence structure as much as possible.

Changes allowed:
- Grammar corrections
- Spelling fixes
- Basic word choice improvements
- Punctuation adjustments

Changes NOT allowed:
- Sentence restructuring
- Changing technical terminology
- Altering the author's voice
"""

# 平衡版 Prompt
"""
You are an academic writing expert. Polish this text to meet
top-tier journal standards while maintaining the core meaning
and important terminology.

Focus on:
- Improving academic vocabulary
- Enhancing sentence fluency
- Strengthening logical connections
- Ensuring grammatical correctness
- Maintaining technical accuracy
"""

# 激进版 Prompt
"""
You are a senior journal editor. Rewrite this text to achieve
publication quality in Nature/Science level journals. Prioritize
clarity, impact, and academic sophistication.

Feel free to:
- Restructure sentences for better flow
- Replace weak phrases with stronger alternatives
- Improve paragraph transitions
- Enhance academic tone significantly
- Optimize for maximum impact
"""
```

### 3.3 修改说明功能（Priority: P1）

#### 功能目标
让用户理解每处修改的价值，提升写作能力

#### 说明维度

```typescript
interface ModificationExplanation {
  // 修改类型
  category: 'vocabulary' | 'grammar' | 'style' | 'structure';

  // 具体问题
  issue: string;  // "使用了口语化表达 'a lot of'"

  // 改进方案
  solution: string;  // "替换为学术表达 'numerous'"

  // 为什么更好
  reason: string;  // "学术写作中应避免非正式表达"

  // 学习建议
  tip: string;  // "常见学术替换：a lot of→numerous/substantial"

  // 相关资源
  reference?: string;  // 链接到写作指南
}
```

#### 常见修改类型库

建立知识库，覆盖常见学术写作问题：

**中式英语类**
- ❌ "more and more" → ✅ "increasingly"
- ❌ "play an important role" → ✅ "contribute significantly to"
- ❌ "do research" → ✅ "conduct research"

**被动语态类**
- 方法部分：多用被动（"The data was analyzed..."）
- 讨论部分：适当主动（"We found that..."）

**时态使用类**
- Introduction: 现在时（"Machine learning is..."）
- Method: 过去时（"We collected data..."）
- Results: 过去时（"The results showed..."）
- Discussion: 现在时（"These findings suggest..."）

### 3.4 质量评估系统（Priority: P1）

#### 评估维度

```go
type QualityAssessment struct {
    // 学术性评分 (0-100)
    AcademicScore int `json:"academic_score"`

    // 可读性评分 (0-100)
    ReadabilityScore int `json:"readability_score"`

    // 详细指标
    Metrics QualityMetrics `json:"metrics"`

    // 问题点
    Issues []Issue `json:"issues"`

    // 改进建议
    Suggestions []string `json:"suggestions"`
}

type QualityMetrics struct {
    // 基础指标
    WordCount           int     `json:"word_count"`
    SentenceCount       int     `json:"sentence_count"`
    AvgSentenceLength   float64 `json:"avg_sentence_length"`    // 建议: 15-25词

    // 学术性指标
    AcademicWordRatio   float64 `json:"academic_word_ratio"`    // 学术词汇占比
    PassiveVoiceRatio   float64 `json:"passive_voice_ratio"`    // 被动语态占比
    NominalRatio        float64 `json:"nominal_ratio"`          // 名词化比例

    // 可读性指标
    FleschReadingEase   float64 `json:"flesch_reading_ease"`    // 60-70为最佳
    ComplexWordRatio    float64 `json:"complex_word_ratio"`

    // 规范性指标
    TenseConsistency    float64 `json:"tense_consistency"`      // 时态一致性
    TerminologyConsistency float64 `json:"terminology_consistency"` // 术语一致性
}

type Issue struct {
    Type        string `json:"type"`        // grammar/style/consistency
    Severity    string `json:"severity"`    // low/medium/high
    Description string `json:"description"`
    Position    int    `json:"position"`
    Suggestion  string `json:"suggestion"`
}
```

#### 评分算法

```go
func CalculateAcademicScore(text string) int {
    score := 50 // 基础分

    // 学术词汇加分 (0-20分)
    academicWordRatio := calculateAcademicWordRatio(text)
    score += int(academicWordRatio * 20)

    // 句式复杂度 (0-15分)
    complexity := calculateSentenceComplexity(text)
    score += complexity

    // 逻辑连接词 (0-10分)
    connectors := countLogicalConnectors(text)
    score += min(connectors * 2, 10)

    // 扣分项
    issues := detectIssues(text)
    for _, issue := range issues {
        if issue.Severity == "high" {
            score -= 5
        } else if issue.Severity == "medium" {
            score -= 2
        }
    }

    return clamp(score, 0, 100)
}
```

---

## 4. 新功能规划

### 4.1 批量处理功能（Priority: P0）

#### 使用场景
用户上传整篇论文（5000-10000词），系统自动分段处理

#### 功能设计

```
┌────────────────────────────────────────────────────┐
│  批量润色                                           │
├────────────────────────────────────────────────────┤
│  📄 已上传: my_paper.docx (6,234 words)            │
│                                                    │
│  ✅ Abstract              [完成] [查看]            │
│  ✅ Introduction          [完成] [查看]            │
│  🔄 Method (处理中... 65%)                         │
│  ⏳ Results                                        │
│  ⏳ Discussion                                     │
│                                                    │
│  进度: 3/5 段落完成                                │
│  预计剩余时间: 2分钟                                │
│                                                    │
│  [暂停] [导出已完成部分] [全部完成后导出]          │
└────────────────────────────────────────────────────┘
```

#### 技术实现

```go
type BatchPolishJob struct {
    JobID       string              `json:"job_id"`
    UserID      int64               `json:"user_id"`
    FileName    string              `json:"file_name"`
    TotalParts  int                 `json:"total_parts"`
    Status      string              `json:"status"`      // pending/processing/completed/failed
    Progress    float64             `json:"progress"`    // 0-100
    CreatedAt   time.Time           `json:"created_at"`
    Parts       []PolishPart        `json:"parts"`
}

type PolishPart struct {
    PartID          string    `json:"part_id"`
    SectionName     string    `json:"section_name"`    // Abstract/Introduction/Method...
    OriginalContent string    `json:"original_content"`
    PolishedContent string    `json:"polished_content"`
    Status          string    `json:"status"`
    ProcessedAt     time.Time `json:"processed_at"`
}

// 批量处理服务
func (s *PolishService) BatchPolish(ctx context.Context, job *BatchPolishJob) error {
    // 1. 保存任务到Redis（用于进度查询）
    s.redis.Set(ctx, "batch:"+job.JobID, job, 24*time.Hour)

    // 2. 异步处理各部分
    for i, part := range job.Parts {
        go func(index int, p PolishPart) {
            // 处理单个段落
            result, err := s.Polish(ctx, &model.PolishRequest{
                Content: p.OriginalContent,
                Style:   "academic",
            }, job.UserID)

            // 更新进度
            job.Parts[index].PolishedContent = result.PolishedContent
            job.Parts[index].Status = "completed"
            job.Progress = float64(index+1) / float64(job.TotalParts) * 100

            // 更新Redis
            s.redis.Set(ctx, "batch:"+job.JobID, job, 24*time.Hour)
        }(i, part)
    }

    return nil
}
```

### 4.2 文件上传与解析（Priority: P0）

#### 支持格式
- `.docx` - Word文档（最常用）
- `.txt` - 纯文本
- `.pdf` - PDF文档（OCR提取）
- `.tex` - LaTeX源文件（未来）

#### 智能分段算法

```go
func ParseDocument(file io.Reader, fileType string) (*Document, error) {
    var doc Document

    switch fileType {
    case "docx":
        doc = parseDocx(file)
    case "txt":
        doc = parseText(file)
    case "pdf":
        doc = parsePDF(file)
    }

    // 智能识别章节
    doc.Sections = identifySections(doc.RawContent)

    // 分段处理（每段500-1000词）
    doc.Parts = splitIntoParts(doc.Sections, 800)

    return &doc, nil
}

type Document struct {
    FileName   string
    FileType   string
    RawContent string
    Sections   []Section   // Abstract, Introduction, Method...
    Parts      []Part      // 分段后的内容块
}

type Section struct {
    Name      string   // 章节名
    Content   string
    StartLine int
    EndLine   int
}

// 章节识别（基于关键词和格式）
func identifySections(content string) []Section {
    sectionKeywords := []string{
        "Abstract", "Introduction", "Related Work", "Background",
        "Method", "Methodology", "Approach",
        "Experiment", "Results", "Evaluation",
        "Discussion", "Conclusion", "References",
    }

    // 使用正则表达式和关键词匹配识别章节
    // ...
}
```

### 4.3 导出功能（Priority: P1）

#### 导出格式

**1. Word修订模式（Track Changes）**
```
最有价值的格式！导师可以直接看到所有修改
- 删除的文字显示为删除线
- 新增的文字显示为下划线/高亮
- 可以接受/拒绝每处修改
```

**2. 对比PDF**
```
左右对照排版，适合打印审阅
- 左侧：原文
- 右侧：润色后
- 修改处标注序号，页面底部有说明
```

**3. 纯净版本**
```
只有润色后的内容，无标记
适合直接提交使用
```

#### 技术实现

```go
// 导出服务
type ExportService struct {
    docxGenerator *DocxGenerator
    pdfGenerator  *PDFGenerator
}

// 导出为Word Track Changes格式
func (e *ExportService) ExportToWordWithChanges(record *PolishRecord) ([]byte, error) {
    doc := docx.New()

    // 计算差异
    diffs := calculateDiff(record.OriginalContent, record.PolishedContent)

    // 添加内容并标记修改
    for _, diff := range diffs {
        switch diff.Type {
        case "delete":
            doc.AddStrikethrough(diff.Text)
        case "insert":
            doc.AddUnderline(diff.Text)
        case "unchanged":
            doc.AddText(diff.Text)
        }
    }

    return doc.Generate()
}
```

### 4.4 学科定制功能（Priority: P2）

#### 设计目标
不同学科有不同的写作规范和术语体系

#### 学科模板

```go
type DisciplineTemplate struct {
    ID          string
    Name        string
    Category    string              // STEM/Medical/Humanities
    Description string

    // 术语库
    Terminology map[string]string   // 标准术语映射

    // 常用句式
    CommonPhrases []Phrase

    // 写作规范
    WritingRules []Rule

    // Prompt增强
    PromptEnhancement string
}

// 示例：计算机科学模板
var ComputerScienceTemplate = DisciplineTemplate{
    ID:   "cs",
    Name: "Computer Science",
    Category: "STEM",

    Terminology: map[string]string{
        "神经网络":    "neural network",
        "深度学习":    "deep learning",
        "注意力机制":  "attention mechanism",
        // ... 更多术语
    },

    CommonPhrases: []Phrase{
        {
            Pattern: "propose_method",
            Examples: [
                "We propose a novel approach that...",
                "This paper introduces a new method for...",
                "We present an innovative technique to...",
            ],
        },
    },

    WritingRules: []Rule{
        {
            Type: "terminology_consistency",
            Description: "确保专业术语前后一致",
        },
        {
            Type: "algorithm_description",
            Description: "算法描述应清晰、可复现",
        },
    },

    PromptEnhancement: `
        This is a Computer Science research paper. Pay special attention to:
        - Use standard CS terminology (e.g., "neural network" not "neural net")
        - Algorithm descriptions should be precise and reproducible
        - Maintain consistent notation for mathematical symbols
        - Follow IEEE/ACM writing conventions
    `,
}
```

#### 使用流程

```
用户选择学科
    ↓
系统加载对应模板
    ↓
Prompt中注入学科特定规则
    ↓
调用AI时使用增强的Prompt
    ↓
后处理：术语一致性检查
```

### 4.5 个人写作档案（Priority: P2）

#### 功能目标
分析用户的写作习惯，提供个性化改进建议

#### 数据收集

```go
type WritingProfile struct {
    UserID        int64
    TotalPolished int        // 总润色次数
    TotalWords    int        // 总字数

    // 常见问题
    CommonIssues  map[string]int  // {"中式英语": 45, "被动语态过度": 23}

    // 进步曲线
    ProgressHistory []ProgressPoint

    // 写作习惯
    WritingHabits WritingHabits

    // 个性化建议
    Recommendations []string
}

type ProgressPoint struct {
    Date          time.Time
    AcademicScore int
    IssuesCount   int
}

type WritingHabits struct {
    AvgSentenceLength  float64  // 平均句长
    PassiveVoiceRatio  float64  // 被动语态比例
    VocabularyDiversity float64 // 词汇多样性
    CommonMistakes     []string // 常犯错误
}
```

#### 分析报告界面

```
┌─────────────────────────────────────────────────────┐
│  📊 你的写作档案                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  总计润色: 127次 | 总字数: 52,341                   │
│  学术性进步: 58分 → 76分 (+18分) 📈                 │
│                                                     │
│  你的写作特点:                                      │
│  ✅ 逻辑清晰，结构合理                              │
│  ⚠️  句子偏长（平均28词，建议20-25词）              │
│  ⚠️  被动语态使用过多（38%，建议20-30%）            │
│                                                     │
│  常见问题 Top 3:                                    │
│  1. 中式英语表达 (23次)                             │
│     常见: "play an important role"                 │
│     建议改为: "significantly contribute to"        │
│                                                     │
│  2. 冠词误用 (18次)                                 │
│     学习资源: [冠词使用指南]                        │
│                                                     │
│  3. 时态不一致 (12次)                               │
│     建议: Method部分统一使用过去时                  │
│                                                     │
│  [查看详细报告] [获取学习计划]                      │
└─────────────────────────────────────────────────────┘
```

### 4.6 引文与术语保护（Priority: P2）

#### 问题场景
AI润色时可能误改：
- 引用格式：`[1]`, `(Zhang et al., 2023)`
- 专业术语：`COVID-19`, `BERT`, `ResNet-50`
- 公式编号：`Equation (1)`, `Table 2`
- 代码片段：`function_name()`

#### 保护机制

```go
// 保护规则定义
type ProtectionRule struct {
    Type    string   // citation/term/formula/code
    Pattern string   // 正则表达式
    Example string
}

var DefaultProtectionRules = []ProtectionRule{
    {
        Type:    "citation",
        Pattern: `\[\d+\]|\([A-Z][a-z]+\s+et\s+al\.,\s+\d{4}\)`,
        Example: "[1] or (Zhang et al., 2023)",
    },
    {
        Type:    "formula",
        Pattern: `(Equation|Formula|Table|Figure)\s+\(?\d+\)?`,
        Example: "Equation (1), Table 2",
    },
    {
        Type:    "code",
        Pattern: `\b[a-z_]+\([^\)]*\)|\b[A-Z][a-z]+\.[a-z_]+`,
        Example: "function_name() or Class.method",
    },
}

// 预处理：替换为占位符
func PreprocessWithProtection(content string, rules []ProtectionRule) (string, map[string]string) {
    protected := make(map[string]string)
    result := content

    for i, rule := range rules {
        re := regexp.MustCompile(rule.Pattern)
        matches := re.FindAllString(content, -1)

        for j, match := range matches {
            placeholder := fmt.Sprintf("__PROTECTED_%s_%d_%d__", rule.Type, i, j)
            protected[placeholder] = match
            result = strings.Replace(result, match, placeholder, 1)
        }
    }

    return result, protected
}

// 后处理：还原
func PostprocessWithProtection(content string, protected map[string]string) string {
    result := content
    for placeholder, original := range protected {
        result = strings.Replace(result, placeholder, original, 1)
    }
    return result
}
```

#### 用户自定义术语库

```
┌─────────────────────────────────────────┐
│  术语保护                                │
├─────────────────────────────────────────┤
│  自动保护:                              │
│  ✅ 引用格式 [1], (Author, 2023)        │
│  ✅ 公式编号 Equation (1)               │
│  ✅ 常见缩写 DNA, RNA, API              │
│                                         │
│  自定义术语:                            │
│  [+ 添加术语]                           │
│                                         │
│  • BERT-base                            │
│  • GPT-3.5-turbo                        │
│  • ResNet-50                            │
│  • your-specific-term                  │
│                                         │
│  [导入术语表] [导出]                    │
└─────────────────────────────────────────┘
```

---

## 5. 技术架构优化

### 5.1 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                        用户层                            │
│  Web界面  │  浏览器插件  │  移动端  │  API接口          │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                      应用层 (Go)                         │
│  ┌──────────┬──────────┬──────────┬──────────┐         │
│  │ 认证服务  │ 润色服务  │ 批处理   │ 导出服务  │         │
│  └──────────┴──────────┴──────────┴──────────┘         │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                      AI服务层                            │
│  ┌──────────────────┬──────────────────┐               │
│  │  Claude API      │  GPT-4 API       │               │
│  │  (主要)          │  (备用)          │               │
│  └──────────────────┴──────────────────┘               │
│  ┌─────────────────────────────────────┐               │
│  │  Prompt管理器                        │               │
│  │  - 学科模板                          │               │
│  │  - 风格控制                          │               │
│  │  - Few-shot案例库                    │               │
│  └─────────────────────────────────────┘               │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│                      数据层                              │
│  ┌──────────┬──────────┬──────────┬──────────┐         │
│  │PostgreSQL│  Redis   │  MinIO   │  ES      │         │
│  │(业务数据) │(缓存/队列)│(文件存储)│(日志/搜索)│         │
│  └──────────┴──────────┴──────────┴──────────┘         │
└─────────────────────────────────────────────────────────┘
```

### 5.2 Prompt工程优化

#### Prompt版本管理

```go
type PromptTemplate struct {
    ID            string
    Version       string
    Name          string
    SystemPrompt  string
    UserPrompt    string
    FewShotExamples []Example
    Parameters    map[string]interface{}
    EffectivenessScore float64  // A/B测试得分
}

type Example struct {
    Input  string
    Output string
    Explanation string
}

// Prompt管理器
type PromptManager struct {
    templates map[string]*PromptTemplate
    db        *sql.DB
}

func (pm *PromptManager) GetBestPrompt(discipline string, style string) *PromptTemplate {
    // 根据学科和风格返回效果最好的Prompt
    key := fmt.Sprintf("%s_%s", discipline, style)
    return pm.templates[key]
}
```

#### 高质量Prompt示例

```
System Prompt:
You are an expert academic editor specializing in {discipline} research papers.
Your goal is to polish academic writing to meet top-tier journal standards
(Nature, Science, Cell, etc.) while preserving the original meaning and
technical accuracy.

Guidelines:
1. PRESERVE: Technical terms, citations, formulas, and core arguments
2. IMPROVE: Vocabulary, sentence structure, clarity, and academic tone
3. ENSURE: Grammatical correctness and logical flow
4. MAINTAIN: Author's voice and writing style (unless style=aggressive)

Current task: Polish the following text with {style} approach
- Conservative: Only fix clear errors
- Balanced: Moderate improvements (RECOMMENDED)
- Aggressive: Significant rewriting for maximum impact

User Prompt:
Please polish the following academic text:

[ORIGINAL TEXT]
{content}
[/ORIGINAL TEXT]

Requirements:
- Target language: {language}
- Writing style: {style}
- Discipline: {discipline}

Output format:
1. Polished version
2. Brief explanation of major changes (2-3 sentences)
3. Confidence score (0-100) for this polishing

---

Few-shot Examples:

Example 1 (Computer Science):
Input: "In this paper, we propose a new method to solve the problem of
image classification using deep learning."

Output: "This paper presents a novel approach to address image
classification challenges through deep learning methodologies."

Explanation: Replaced informal "solve the problem" with academic
"address challenges"; "new method" → "novel approach"; added
discipline-appropriate terminology.

Example 2 (Medical Science):
Input: "The results show that the drug works well in treating the disease."

Output: "The results demonstrate that the therapeutic agent exhibits
significant efficacy in disease management."

Explanation: Elevated vocabulary: "show" → "demonstrate", "drug" →
"therapeutic agent", "works well" → "exhibits significant efficacy",
"treating" → "management".
```

### 5.3 性能优化

#### 缓存策略

```go
type CacheService struct {
    redis *redis.Client
}

// 缓存key生成
func (c *CacheService) GenerateCacheKey(req *PolishRequest) string {
    // 对内容+参数进行hash
    data := fmt.Sprintf("%s:%s:%s:%s",
        req.Content, req.Style, req.Language, req.Provider)
    hash := sha256.Sum256([]byte(data))
    return fmt.Sprintf("polish:%x", hash)
}

// 查询缓存
func (c *CacheService) GetCached(req *PolishRequest) (*PolishResponse, error) {
    key := c.GenerateCacheKey(req)
    data, err := c.redis.Get(context.Background(), key).Result()
    if err != nil {
        return nil, err
    }

    var resp PolishResponse
    json.Unmarshal([]byte(data), &resp)
    return &resp, nil
}

// 设置缓存（7天过期）
func (c *CacheService) SetCache(req *PolishRequest, resp *PolishResponse) {
    key := c.GenerateCacheKey(req)
    data, _ := json.Marshal(resp)
    c.redis.Set(context.Background(), key, data, 7*24*time.Hour)
}
```

**缓存命中率预估**：
- 相同段落重复润色：20-30%
- 示例文本和测试：10-15%
- 预期总命中率：30-45%

#### 批量处理优化

```go
// 使用worker pool并发处理
type PolishWorkerPool struct {
    workers   int
    jobs      chan PolishJob
    results   chan PolishResult
    semaphore chan struct{}
}

func NewPolishWorkerPool(workers int) *PolishWorkerPool {
    return &PolishWorkerPool{
        workers:   workers,
        jobs:      make(chan PolishJob, workers*2),
        results:   make(chan PolishResult, workers*2),
        semaphore: make(chan struct{}, workers),
    }
}

func (p *PolishWorkerPool) ProcessBatch(jobs []PolishJob) []PolishResult {
    // 启动workers
    for i := 0; i < p.workers; i++ {
        go p.worker()
    }

    // 提交任务
    go func() {
        for _, job := range jobs {
            p.jobs <- job
        }
        close(p.jobs)
    }()

    // 收集结果
    var results []PolishResult
    for i := 0; i < len(jobs); i++ {
        results = append(results, <-p.results)
    }

    return results
}
```

### 5.4 成本控制

#### AI调用成本估算

**Claude API定价（2024年）**:
- Input: $0.003 / 1K tokens
- Output: $0.015 / 1K tokens

**典型场景**:
```
单次润色（500词）:
- Input: 500词 ≈ 650 tokens → $0.00195
- Output: 500词 + 说明 ≈ 700 tokens → $0.0105
- 单次成本: ≈ $0.012 (约￥0.09)

每日1000次润色:
- 成本: $12/天 ≈ ￥360/月
- 收入: 200用户 × ￥20/月 = ￥4000/月
- 毛利率: 91%
```

#### 成本优化策略

```go
// 1. 智能路由：简单任务用小模型
func (s *PolishService) SelectModel(content string) string {
    wordCount := len(strings.Fields(content))

    if wordCount < 100 {
        return "claude-3-haiku"  // 更便宜
    } else if wordCount < 500 {
        return "claude-3-sonnet"  // 平衡
    } else {
        return "claude-3-opus"    // 高质量
    }
}

// 2. 流式输出：减少等待时间感知
func (s *PolishService) PolishWithStreaming(ctx context.Context, req *PolishRequest) <-chan string {
    stream := make(chan string)

    go func() {
        defer close(stream)
        // 流式调用AI API
        // 逐步返回结果
    }()

    return stream
}

// 3. 批量折扣：合并短请求
func (s *PolishService) BatchSmallRequests(requests []*PolishRequest) {
    // 将多个小段落合并成一个请求
    // 减少API调用次数
}
```

---

## 6. 商业模式设计

### 6.1 定价策略

#### 免费层（Free Tier）
**目标**: 获客、让用户体验核心价值

```
✅ 每天 3 次免费润色
✅ 单次最多 300 词
✅ 基础风格（academic）
✅ 对比展示
❌ 无历史记录
❌ 无批量处理
❌ 无导出功能
```

#### 标准版（Standard）- ￥29/月 或 ￥199/年
**目标**: 个人硕博士用户

```
✅ 无限次数润色
✅ 单次最多 1000 词
✅ 3种风格 + 3个备选版本
✅ 完整对比和修改说明
✅ 历史记录（保存30天）
✅ 批量处理（最多10段）
✅ 导出Word/PDF
✅ 邮件支持
```

#### 专业版（Pro）- ￥79/月 或 ￥699/年
**目标**: 高频用户、即将投稿用户

```
✅ Standard 的所有功能
✅ 单次最多 3000 词
✅ 学科定制（10+学科模板）
✅ 历史记录永久保存
✅ 批量处理（无限制）
✅ 个人写作档案和进步分析
✅ 术语库和引文保护
✅ 优先使用最新AI模型
✅ 优先技术支持
```

#### 团队版（Team）- ￥299/月（5人）
**目标**: 实验室、研究组

```
✅ Pro 的所有功能（每人）
✅ 团队协作和分享
✅ 统一管理和账单
✅ 自定义学科模板
✅ API接口访问
✅ 专属客户经理
✅ 发票支持
```

### 6.2 收入预测

#### 乐观场景（12个月）

```
月份 | 免费用户 | 付费用户 | 月收入 | 累计收入
-----|---------|---------|--------|----------
M1   | 100     | 5       | ￥145  | ￥145
M2   | 300     | 15      | ￥435  | ￥580
M3   | 500     | 30      | ￥870  | ￥1,450
M4   | 800     | 60      | ￥1,740| ￥3,190
M5   | 1,200   | 100     | ￥2,900| ￥6,090
M6   | 1,800   | 180     | ￥5,220| ￥11,310
M7   | 2,500   | 280     | ￥8,120| ￥19,430
M8   | 3,500   | 420     | ￥12,180| ￥31,610
M9   | 5,000   | 600     | ￥17,400| ￥49,010
M10  | 7,000   | 850     | ￥24,650| ￥73,660
M11  | 9,000   | 1,150   | ￥33,350| ￥107,010
M12  | 12,000  | 1,500   | ￥43,500| ￥150,510

假设:
- 转化率: 5-12% (免费→付费)
- 平均客单价: ￥29/月
- 月增长率: 30-40%
```

### 6.3 推广策略

#### 初期（0-3个月）：种子用户

**目标**: 100个真实用户，10个付费

**策略**:
1. **内容营销**
   - 知乎：写5篇"学术论文写作"相关回答
   - 小红书：发布"论文润色前后对比"案例
   - B站：录制"AI论文润色教程"视频

2. **社群运营**
   - 加入10个硕博QQ群/微信群
   - 提供免费试用码（1个月Pro）
   - 收集反馈，快速迭代

3. **导师推荐计划**
   - 邀请5位导师试用
   - 导师推荐学生，双方获益
   - 建立信任和口碑

#### 中期（3-9个月）：规模化增长

**目标**: 5000用户，500付费

**策略**:
1. **SEO优化**
   - 关键词："论文润色"、"学术英文润色"、"Paper Polish AI"
   - 内容矩阵：写作技巧博客

2. **合作推广**
   - 与科研工具（Zotero、Overleaf）合作
   - 学术社区（科研之友、小木虫）投放
   - 高校社团合作

3. **转介绍激励**
   - 推荐好友，双方各得1周Pro
   - 成功推荐3人，获得1个月免费

#### 后期（9-24个月）：品牌建立

**目标**: 50000用户，5000付费，月收入￥15万

**策略**:
1. **品牌建设**
   - 成为"学术论文润色"领域的第一选择
   - 发布年度"学术写作质量报告"
   - 举办线上写作工作坊

2. **企业服务**
   - 推出高校/科研机构企业版
   - API服务，集成到其他工具

3. **国际化**
   - 英文界面
   - 支持更多学科和语言

---

## 7. 产品建议

### 7.1 用户体验原则

#### 1. 速度至上
```
用户期望:
- 输入后 < 3秒 看到结果
- 批量处理有清晰进度反馈
- 不要让用户等待太久

实现:
✅ 使用流式输出
✅ 缓存常见内容
✅ 后台预加载
✅ 进度条动画
```

#### 2. 简单易用
```
设计理念:
- 核心功能一键完成
- 高级功能渐进展示
- 避免复杂的配置

实例:
❌ 复杂: 需要选择10个参数才能开始
✅ 简单: 粘贴→点击润色→完成（高级选项隐藏在"更多"中）
```

#### 3. 信任建立
```
用户顾虑: AI会不会乱改？改得对不对？

解决方案:
✅ 每处修改都可以查看理由
✅ 提供多个版本对比
✅ 用户可以拒绝任何修改
✅ 显示AI置信度
✅ 保留历史版本
```

#### 4. 学习价值
```
不只是工具，更是老师:
✅ 解释为什么这样改
✅ 提供写作建议
✅ 跟踪用户进步
✅ 推荐学习资源
```

### 7.2 技术债务管理

#### 需要避免的问题

1. **过早优化**
   - ❌ 不要一开始就搭建复杂的微服务架构
   - ✅ 先用单体应用，功能验证后再拆分

2. **功能堆砌**
   - ❌ 不要一次性开发所有功能
   - ✅ MVP → 核心功能 → 高级功能，逐步迭代

3. **忽视测试**
   - ❌ 不要因为赶进度跳过测试
   - ✅ 核心功能必须有单元测试和集成测试

4. **文档缺失**
   - ❌ 不要认为"代码即文档"
   - ✅ API文档、架构文档、运维文档必须齐全

### 7.3 数据驱动决策

#### 关键指标追踪

```go
type ProductMetrics struct {
    // 用户指标
    DAU             int     // 日活跃用户
    MAU             int     // 月活跃用户
    RetentionRate7D float64 // 7日留存率
    ChurnRate       float64 // 流失率

    // 使用指标
    AvgPolishPerUser  float64 // 人均润色次数
    AvgWordsPerPolish int     // 平均每次字数
    FeatureUsageRate  map[string]float64 // 功能使用率

    // 质量指标
    UserSatisfaction  float64 // 用户满意度（1-5分）
    AcceptanceRate    float64 // 修改接受率
    ErrorRate         float64 // 错误率

    // 商业指标
    ConversionRate    float64 // 转化率（免费→付费）
    MRR               float64 // 月经常性收入
    LTV               float64 // 用户生命周期价值
    CAC               float64 // 获客成本
    LTVCAC            float64 // LTV/CAC比率（健康值>3）
}
```

#### A/B测试框架

```go
type ABTest struct {
    TestID      string
    Feature     string
    VariantA    string  // 对照组
    VariantB    string  // 实验组
    SplitRatio  float64 // 流量分配（0.5 = 50/50）
    Metrics     []string
    StartTime   time.Time
    Status      string
}

// 示例：测试不同Prompt效果
func TestPromptVariants() {
    test := ABTest{
        TestID:   "prompt_v2_vs_v3",
        Feature:  "polish_prompt",
        VariantA: "prompt_v2",
        VariantB: "prompt_v3",
        SplitRatio: 0.5,
        Metrics: []string{"user_satisfaction", "acceptance_rate", "processing_time"},
    }

    // 运行1周后分析结果
    results := AnalyzeABTest(test)
    if results.Winner == "VariantB" {
        // 全量切换到v3
        DeployPrompt("prompt_v3")
    }
}
```

### 7.4 风险管理

#### 潜在风险

1. **AI服务中断**
   ```
   风险: Claude/GPT API故障或限流

   应对:
   ✅ 多提供商热备份
   ✅ 本地缓存提高命中率
   ✅ 降级方案（基础语法检查）
   ✅ 用户友好的错误提示
   ```

2. **成本失控**
   ```
   风险: 用户量暴增导致AI调用费用激增

   应对:
   ✅ 设置单用户限流
   ✅ 实时成本监控和告警
   ✅ 准备融资或应急资金
   ```

3. **学术诚信争议**
   ```
   风险: 被质疑"代写"或违反学术规范

   应对:
   ✅ 明确定位为"辅助工具"
   ✅ 用户协议中声明责任
   ✅ 网站显著位置说明合规使用
   ✅ 与高校合作建立最佳实践
   ```

4. **数据安全**
   ```
   风险: 用户论文内容泄露

   应对:
   ✅ 数据加密存储
   ✅ 严格的权限控制
   ✅ 定期安全审计
   ✅ 提供"阅后即焚"模式（不保存历史）
   ```

---

## 8. 实施路线图

### 8.1 MVP阶段（第1-2个月）

**目标**: 验证核心价值，获得100个真实用户

#### 第1周：基础界面
- [ ] 简洁的Web界面（单页应用）
- [ ] 文本输入框 + 一键润色按钮
- [ ] 显示润色结果（纯文本）
- [ ] 基础的用户注册/登录

**交付物**: 可用的Demo，能展示给用户

#### 第2-3周：对比功能（P0）
- [ ] 左右对照展示界面
- [ ] 高亮差异显示
- [ ] 差异计算算法（diff algorithm）
- [ ] 修改分类（词汇/语法/句式）
- [ ] 简单的统计信息

**交付物**: 核心差异化功能，形成产品壁垒

#### 第4-5周：多版本润色（P0）
- [ ] 实现3种润色强度的Prompt
- [ ] 并行调用AI生成3个版本
- [ ] 版本切换和对比界面
- [ ] 性能优化（缓存、并发控制）

**交付物**: 给用户选择权，提升满意度

#### 第6-8周：用户体验优化
- [ ] 历史记录列表和查询
- [ ] 复制、导出纯文本功能
- [ ] 错误处理和用户提示
- [ ] 响应式设计（移动端适配）
- [ ] Loading状态和进度提示

**交付物**: 完整的MVP，可以邀请用户测试

#### 第8周：上线和推广
- [ ] 部署到生产环境
- [ ] 监控和日志系统
- [ ] 准备推广材料（截图、视频）
- [ ] 邀请50-100个种子用户
- [ ] 收集反馈，建立用户群

**里程碑**: MVP上线，开始真实用户测试

### 8.2 核心功能阶段（第3-5个月）

**目标**: 完善核心功能，实现产品PMF，达到500用户，50付费

#### 第9-11周：批量处理（P0）
- [ ] 文件上传功能（.docx, .txt）
- [ ] 文档解析和分段算法
- [ ] 批量任务队列（Redis）
- [ ] 进度追踪和实时更新
- [ ] 断点续传和失败重试

**价值**: 解决全文润色痛点，提升效率10倍

#### 第12-14周：修改说明（P1）
- [ ] 为每处修改生成说明
- [ ] 修改理由分类（语法/词汇/风格）
- [ ] 悬停提示和交互
- [ ] 学习建议和相关资源链接
- [ ] 常见问题知识库

**价值**: 提供学习价值，增加用户粘性

#### 第15-17周：质量评估（P1）
- [ ] 实现评分算法
- [ ] 多维度指标计算
- [ ] 可视化报告界面
- [ ] 问题检测和高亮
- [ ] 改进建议生成

**价值**: 让用户看到提升，增强信任感

#### 第18-20周：付费系统
- [ ] 定价页面设计
- [ ] 支付接口集成（微信/支付宝）
- [ ] 订阅管理系统
- [ ] 使用限额控制
- [ ] 发票和收据

**里程碑**: 开始商业化，获得第一批付费用户

### 8.3 专业化阶段（第6-9个月）

**目标**: 建立专业壁垒，达到5000用户，500付费，月收入￥1.5万

#### 第21-24周：学科定制（P2）
- [ ] 10个学科模板开发
- [ ] 术语库构建（每个学科500+术语）
- [ ] 学科专用Prompt优化
- [ ] 用户选择和自动识别
- [ ] A/B测试不同模板效果

**价值**: 专业性壁垒，竞争对手难以复制

#### 第25-28周：导出功能（P1）
- [ ] Word导出（docx格式）
- [ ] Word Track Changes格式
- [ ] PDF对比导出
- [ ] LaTeX支持（基础）
- [ ] 模板和样式自定义

**价值**: 融入用户工作流，提升留存

#### 第29-32周：引文保护（P2）
- [ ] 保护规则引擎
- [ ] 自动识别引用和术语
- [ ] 用户自定义术语库
- [ ] 术语一致性检查
- [ ] 批量导入导出术语

**价值**: 确保学术规范，建立专业信任

#### 第33-36周：写作档案（P2）
- [ ] 用户数据分析模块
- [ ] 常见问题识别算法
- [ ] 进步曲线可视化
- [ ] 个性化建议生成
- [ ] 学习计划推荐

**里程碑**: 产品成熟，建立专业品牌

### 8.4 规模化阶段（第10-12个月）

**目标**: 扩大市场份额，达到20000用户，2000付费，月收入￥6万

#### 第37-40周：性能优化
- [ ] 系统性能压测
- [ ] 数据库查询优化
- [ ] 缓存策略优化
- [ ] CDN加速（静态资源）
- [ ] 微服务拆分（如需要）

#### 第41-44周：团队协作功能
- [ ] 团队账号和权限管理
- [ ] 分享和协作
- [ ] 评论和反馈
- [ ] 团队统计和报表
- [ ] 企业版定价

#### 第45-48周：移动端和插件
- [ ] 移动端H5优化
- [ ] Chrome浏览器插件
- [ ] 微信小程序（可选）
- [ ] API文档和SDK
- [ ] 第三方集成（Overleaf, Zotero）

#### 第49-52周：品牌和内容
- [ ] 官网SEO优化
- [ ] 内容营销矩阵
- [ ] 线上活动和工作坊
- [ ] 用户案例和Success Story
- [ ] 年度报告和白皮书

**里程碑**: 成为细分领域知名品牌

---

## 9. 成功指标

### 9.1 产品指标（Product Metrics）

#### 用户增长
```
目标值（12个月）:
- 注册用户: 20,000+
- 月活用户: 8,000+
- 日活用户: 2,000+
- 7日留存: 40%+
- 30日留存: 25%+
```

#### 用户参与度
```
目标值:
- 平均每用户润色次数: 15次/月
- 平均每次字数: 500词
- 批量处理使用率: 30%
- 历史记录查看率: 50%
- 功能使用广度: 平均使用5个功能
```

#### 质量指标
```
目标值:
- 用户满意度: 4.5/5.0
- 修改接受率: 85%+
- 错误率: < 1%
- 响应时间: < 5秒（P95）
- 系统可用性: 99.5%
```

### 9.2 商业指标（Business Metrics）

#### 转化和收入
```
目标值（12个月）:
- 免费→付费转化率: 10%+
- 付费用户数: 2,000
- 月经常性收入(MRR): ￥60,000
- 年度经常性收入(ARR): ￥720,000
- 客单价: ￥30/月
```

#### 用户价值
```
目标值:
- 用户生命周期价值(LTV): ￥360
- 获客成本(CAC): ￥50
- LTV/CAC比率: 7.2 (健康值>3)
- 月流失率: < 5%
- 年续费率: 70%+
```

#### 成本结构
```
预期（月收入￥60,000）:
- AI API成本: ￥6,000 (10%)
- 服务器成本: ￥3,000 (5%)
- 营销成本: ￥15,000 (25%)
- 人力成本: ￥20,000 (33%)
- 毛利率: 27% → 净利润: ￥16,000/月
```

### 9.3 里程碑检查点

#### MVP阶段（2个月）
- ✅ 产品上线，功能可用
- ✅ 100个真实用户使用
- ✅ 用户反馈积极（满意度>4.0）
- ✅ 核心功能验证（对比、多版本）

**决策点**: 如果达成，继续；如果未达成，分析原因，调整方向

#### PMF阶段（5个月）
- ✅ 500用户，其中50付费
- ✅ 用户主动推荐（NPS>30）
- ✅ 留存率达标（7日>40%）
- ✅ 有用户愿意长期付费

**决策点**: 达成PMF，加大投入；未达成，重新思考产品定位

#### 规模化阶段（9个月）
- ✅ 5,000用户，500付费
- ✅ 月收入￥15,000
- ✅ 正向现金流
- ✅ 形成品牌认知度

**决策点**: 考虑融资或自力更生继续增长

#### 成熟阶段（12个月）
- ✅ 20,000用户，2,000付费
- ✅ 月收入￥60,000
- ✅ 稳定盈利
- ✅ 行业领先地位

**决策点**: 考虑国际化、多元化或被收购

---

## 10. 附录

### 10.1 技术栈选择

#### 后端
- **语言**: Go 1.21+ (已选择)
- **框架**: Gin (已使用)
- **数据库**: PostgreSQL 15+ (已使用)
- **缓存**: Redis 7+
- **消息队列**: Redis Streams or RabbitMQ
- **文件存储**: MinIO or 阿里云OSS

#### 前端
- **框架**: React 18 + TypeScript
- **UI库**: Ant Design or Chakra UI
- **状态管理**: Zustand or Redux Toolkit
- **构建工具**: Vite

#### AI服务
- **主要**: Claude 3.5 Sonnet
- **备用**: GPT-4 Turbo
- **未来**: 国产大模型（通义千问、文心一言）

#### 运维
- **容器化**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **监控**: Prometheus + Grafana
- **日志**: ELK Stack (Elasticsearch + Logstash + Kibana)
- **错误追踪**: Sentry

### 10.2 竞品分析

#### Grammarly
- **优势**: 品牌知名度高、功能全面、实时检查
- **劣势**: 不够学术化、定价高（$12-30/月）、通用英文为主
- **差异化**: 我们更专注学术、更理解学科术语、性价比更高

#### QuillBot
- **优势**: 改写能力强、价格亲民（$9.95/月）
- **劣势**: 学术性不足、无学科定制、缺乏深度分析
- **差异化**: 我们提供修改说明、学科模板、写作档案

#### Trinka AI
- **优势**: 专注学术写作、理解学科差异
- **劣势**: 界面老旧、功能单一、知名度低
- **差异化**: 更现代的产品体验、批量处理、学习功能

#### DeepL Write
- **优势**: 德系品质、改写自然
- **劣势**: 缺乏学术特色、功能简单
- **差异化**: 学科定制、质量评估、中文市场

### 10.3 学习资源

#### 学术写作指南
- *The Elements of Style* by Strunk & White
- *Academic Writing: A Handbook for International Students*
- Nature/Science 投稿指南
- IEEE Author Center

#### AI Prompt工程
- OpenAI Prompt Engineering Guide
- Anthropic Claude Best Practices
- Prompt Engineering Institute

#### 产品设计
- *Inspired* by Marty Cagan
- *The Lean Startup* by Eric Ries
- *Hooked* by Nir Eyal

---

## 总结

这份产品规划文档提出了一个**专注、深耕学术论文润色**的产品战略。核心理念是：

1. **做一件事，做到极致** - 放弃平台模式，聚焦垂直领域
2. **建立专业壁垒** - 学科定制、质量评估、写作档案
3. **提供学习价值** - 不只是改，更要让用户知道为什么
4. **渐进式增长** - MVP验证 → PMF → 规模化 → 品牌化

**关键成功要素**:
- ✅ 极致的对比展示体验
- ✅ 多版本润色给用户选择权
- ✅ 修改说明提供学习价值
- ✅ 学科定制建立专业壁垒
- ✅ 合理定价和清晰的增长路径

**执行建议**:
1. 先完成MVP，快速上线验证
2. 找到100个种子用户，深度访谈
3. 根据反馈快速迭代，追求PMF
4. 达到PMF后，加大投入规模化
5. 始终保持产品质量和用户体验

祝产品成功！🚀

---

**文档维护**:
- 定期更新（建议每月）
- 根据实际进展调整计划
- 记录决策和经验教训
- 与团队共享，保持一致
