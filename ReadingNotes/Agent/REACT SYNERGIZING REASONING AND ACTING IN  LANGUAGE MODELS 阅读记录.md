# REACT: SYNERGIZING REASONING AND ACTING IN  LANGUAGE MODELS 阅读记录

原文链接：[ReAct: Synergizing Reasoning and Acting in Language Models](http://arxiv.org/abs/2210.03629)

## Intro

传统的Agent通常将 **Reasoning（推理）** 和 **Acting（行动）**分开：

1. **Reasoning-only（如 CoT）**：能够进行复杂推理，但主要依赖模型内部知识，无法通过行动从外部环境获取信息，容易产生幻觉并导致错误传播。*（CoT：chain-of-thought）*
2. **Acting-only**：能够与环境交互，但缺乏推理和规划能力，面对复杂任务时容易产生错误行动。

本文提出了 **ReAct** ：将 Reasoning 和 Acting 结合起来，让 LLM 在任务执行过程中交错生成 **Thought（推理）** 和 **Action（行动）**，并利用 Action 得到的 **Observation（环境反馈）** 指导后续推理和行动。

## ReAct

**Overview**

![img](https://cdn.nlark.com/yuque/0/2026/png/50775921/1790431284099-9fd04a06-3f4f-4fe5-b2d8-09d5821ff5d9.png?x-oss-process=image%2Fformat%2Cwebp)

**Key**

ReAct 的思想很简单 本质将 Agent 的 action space 从 \(\mathcal A\) 扩展为 \( \hat{\mathcal A}=\mathcal A\cup\mathcal L \)，即允许模型除了生成实际 Action 外，还生成自然语言形式的 Thought。    *\(\mathcal{L}\) 表示语言空间（language space）*

当  \(\hat{\mathcal a}  \in  \mathcal L\)  将其称为一个 **thought（思考）** 或 **reasoning trace（推理轨迹）**。

其实就是  Action Space = Action + Language Thought.

**Thought**：不会改变外部环境，也不会产生 Observation，而是基于现在的Context进行推理，从中组织出有用的信息制定/调整行动计划并**更新Context**。

Thought不能只简单的理解为 "思考",在这一过程中还可以承担：

1. **Task Decomposition**：分解任务目标

2. **Action Planning**：制定行动计划
3. **Knowledge Injection** ：引入常识等知识

4. **Information Extraction**：从 Observation 中提取重要信息

5. **Progress Tracking**：跟踪任务进度

6. **Plan Adjustment**：根据环境反馈调整行动计划

7. **Exception Handling**：处理异常情况

## 实验部分

![image-20261004003510296](./REACT SYNERGIZING REASONING AND ACTING IN  LANGUAGE MODELS 阅读记录.assets/image-20261004003510296.png)

这里结合论文给出的具体的例子，阐述一下ReAct优于Reasoning-only和Acting-only。



### LLM Agent实验常用的数据集

#### 验证推理任务

**HotpotQA** ：是一个**多跳问答（Multi-hop Question Answering）数据集**。多跳就是回答一个问题需要连接多个信息来源/多个推理步骤，而不是从一篇文档中直接找到答案。
**例子**：**先**查"Apple Remote 最初控制什么程序"，得到 **Front Row**，再查"什么设备可以控制 Front Row"，最终得到答案。



**FEVER**  ：**Fact Verification（事实验证）**，给模型一个 claim（声明），让模型判断这个声明是否被证据**支持（SUPPORTS）**、**反驳（REFUTES）或无法确定（NOT ENOUGH INFO**）**。**模型不能完全依赖自己的记忆，而是可以主动寻找证据验证事实。

**例子**：“Apple Remote 于 2005 年推出。”→ Agent 搜索相关资料 → 根据证据判断该说法是否成立。



#### 验证决策任务

**ALFWorld**：一个**文本形式的交互式家庭环境**，Agent 需要通过连续的 Action 与环境交互来完成任务。

**例子**：“把一个杯子放进柜子” → 寻找杯子 → 拿起杯子 → 找到柜子 → 打开柜子 → 放入杯子。



**WebShop**：一个**模拟网上购物的交互环境**，Agent 根据用户需求搜索、浏览和选择商品。

**例子**：“寻找一个价格低于 $100、黑色、带抽屉的床头柜” → 搜索商品 → 查看属性 → 筛选 → 选择符合要求的商品。

