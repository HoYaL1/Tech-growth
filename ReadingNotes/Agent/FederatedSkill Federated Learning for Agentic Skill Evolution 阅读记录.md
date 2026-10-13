# FederatedSkill: Federated Learning for Agentic Skill Evolution 阅读记录

原文链接：[FederatedSkill: Federated Learning for Agentic Skill Evolution](https://arxiv.org/abs/2606.03143 "arXiv:2606.03143")



------



## BG

​	现代大语言模型智能体（LLM agents）越来越依赖**技能库**处理复杂任务，使得**技能演化**成为智能体自我提升的核心驱动力。然而，孤立的单用户任务流缺少构建完备技能集所需的数据多样性。跨用户协作虽能突破这一数据瓶颈，但现有的轨迹共享方案会损害用户隐私，且采用一刀切的全局技能库，无法适配异构客户端的差异化需求。

​	本文提出 FederatedSkill：一个面向协作式智能体演化的隐私保护框架。FederatedSkill 摒弃原始轨迹共享的思路，将**语义技能差异（semantic skill diffs）**—— 即作用于本地技能库的结构化补丁 —— 作为**通信的基本单元**。在服务端，一个演化智能体对这些补丁进行聚合，动态建模各个客户端独有的能力边界，以此实现高度个性化的技能演化，而非生成效果欠佳的全局平均模型。

## 问题

C1：直接上传**原始轨迹（Raw Trajectories）**会导致隐私泄露

C2：一个**全局技能库（Global Skill Library）**无法满足个性化需求

## 解决方案

R1：只上传Skill Patch

R2：专门设置一个Evolution Agent（一个由 Skill 驱动的 Agent）



------



## FederatedSkill Detailed Information

### Client-Side

#### Trial Execution（试验执行）

客户端 Agent 使用当前技能库执行本地任务，收集**任务执行轨迹**。

[^trajectory]: *<span style="font-size: 0.75em;">**任务执行轨迹(trajectory)**在大部分情况下本质就是**短期记忆**，本文则是智能体<strong>单次任务的完整行动链</strong>：思考、工具调用、环境反馈、结果，是<strong>任务执行过程产生的原始事件记录</strong>。</span>*

在第 \(t\) 轮联邦训练中，第 \(i\) 个客户端从本地任务分布 \(D_i\) 中采样一组任务，并使用当前技能库 \(L_i^t\) 执行任务。Agent 执行任务后生成执行轨迹 \(\tau\)，环境根据任务完成情况对轨迹进行验证并给出**奖励信号 \(R_{i,x}(\tau)\)**。

#### Per-Client Patch Distillation（客户端个性化补丁蒸馏）

每个客户端设置一个本地 **Patcher（补丁生成器）**，它使用与执行 Agent 相同的 LLM Backbone，将本地试验产生的原始记录转化为结构化 Skill Patch。

##### A.Patch Generation（补丁生成）

论文将**补丁生成**设计为一次 **LLM 调用**，并向 Patcher 提供三类输入。

1. 输入一：压缩后的执行轨迹（Compacted Trajectory）
   原始轨迹经过压缩后，只保留最多 \(K_{\text{step}}\) 个 Agent 执行步骤，即初始步骤与最近的 \(K_{\text{step}}-1\) 个步骤。执行指标会被移除，环境 Observation 则截断至 \(K_{\text{obs}}\) 个字符，并用明确的 `<truncated>` 标记提示内容已被截断。
   **目的：**控制上下文长度与计算成本，同时保留对技能改进有价值的执行信息。

   [^Compacted Trajectory]: <span style="font-size: 0.75em;">轨迹压缩不等于直接把原始轨迹上传到服务器。压缩后的轨迹仅作为客户端本地蒸馏过程的输入，最终对外传输的是生成后的 Skill Patch。</span>

2. 输入二：当前技能库快照（Library Snapshot）
   将客户端当前技能库 \(L_i^t\) 表示为 JSON 快照，其中包含各个文件的路径和对应内容。
   **目的：**让 Patcher 知道客户端已经拥有哪些技能，避免在缺乏现有上下文的情况下盲目创建或修改技能。

3. 输入三：试验结果（Trial Outcome）
   包含任务名称、质量或奖励信号 \(R_{i,x}(\tau)\)、异常类型、Agent 最终输出，以及可用的验证子测试失败信息。
   **目的：**帮助 Patcher 区分执行成功与失败的情况，并结合失败原因判断技能应当如何修改。

##### B. Patch Schema（补丁结构）

论文将客户端**生成的补丁**定义为一个四元组：

​									\[ \delta_i^t = \left(U_i^t,D_i^t,R_{i,x}(\tau),s_i^t\right) \]

<div style="font-size: 0.85em;">

| 符号              | 专业术语                 | 含义                                       |
| :---------------- | :----------------------- | :----------------------------------------- |
| \(U_i^t\)         | File Upserts             | 新增或更新的文件及其完整内容               |
| \(D_i^t\)         | File Deletions           | 建议从当前技能库删除的文件路径             |
| \(R_{i,x}(\tau)\) | Trial Reward             | 本次任务执行的验证奖励信号                 |
| \(s_i^t\)         | Natural-language Summary | 用一句自然语言概括本次技能修改的原因或内容 |

##### C.隐私保护设计

**补丁 \(\delta_i^t\) 是唯一传出客户端的信息载体**。**原始**试执行**数据**全部**保留在本地**节点，包括完整轨迹、验证器输出、中间文件以及任务输入。

### Server-Side

**核心思想**：***服务器不把所有客户端上传的 Skill Patch 简单合并成一个统一技能库，而是分析不同客户端的能力差异，再分别演化每个客户端的技能库*。**

在每一轮结束时，服务器收集来自 \(N\) 个客户端、针对同一任务族（Task Family）的一组技能补丁，随后，服务器分别对每个客户端的技能库进行演化。  Evolution Agent \(M\) 在**部分可观测（Partially Observable）**的条件下运行。它能够观察到backbone model 和agent harness 以及skill patch。但无法访问原始任务执行轨迹（Raw Trajectories），因此必须仅根据技能补丁的语义推断每个客户端的能力边界（Capability Boundaries）。

#### Stage 1: Evolution Planning

在本阶段Evolution Agent将本轮的技能补丁提炼为演化计划 \(P_t\)，其中包含三个组成部分：

1. **能力矩阵（Capability Matrix）：**定位每个客户端在不同任务类型上的能力状态。
2. **双层能力记忆（Two-Level Capability Memory）：**记录客户端的能力边界及其失败原因。
3. **演化指令（Evolution Directives）：**规定下一步应该采取哪些技能新增、修复或重构操作。

#### Stage 2: Per-Client Library Evolution  Executing the Plan.

本阶段对每个客户端独立运行，以执行演化计划 \(P_t\) 中的各项指令。

对于客户端 \(i\)，Evolution Agent 会综合处理以下信息：

- 当前技能库 \(L_i^t\)；
- 本轮收到的技能补丁；
- 其他客户端技能库的**只读快照（Read-only Snapshot）**。

随后，Evolution  Agent 根据演化指令修改客户端 \(i\) 的技能库 \(L_i^t\)，具体操作包括：

1. **吸收可迁移补丁（Absorbing Transferable Patches）**：将其他客户端的有效技能经验应用到当前客户端。
2. **原位修复失效技能（Repairing Broken Skills in Place）**：直接修复当前技能库中存在但无法正常工作的技能。
3. **重写其他客户端的技能（Rewriting Peer Skills）**：根据客户端 \(i\) 的特征 \(\rho_i\)，改写其他客户端的技能，使其适配当前客户端的基础模



------



## 一些思考

| 传统联邦学习                  | FederatedSkill                              |
| ----------------------------- | ------------------------------------------- |
| 本地训练数据 \(D_i\)          | 本地任务轨迹 \(\mathcal B_i^t\)             |
| 本地模型参数 \(\theta_i^t\)   | 当前技能库 \(\mathcal L_i^t\)               |
| 参数更新 \(\Delta\theta_i^t\) | 技能补丁 \(\delta_i^t\)                     |
| 聚合后的模型                  | 个性化更新后的技能库 \(\mathcal L_i^{t+1}\) |

如果按上面的理解对标传统的联邦学习的话，那么在技能库这里就出现了问题（个人理解），原文赋予EAgent的权利是可以看其他的client的a read-only snapshot of peer libraries，论文对**Library Snapshot**的定义: A JSON snapshot of the current library  \(\mathcal L_i^t\)detailing all file paths and **their corresponding contents**.

Skill 里本身可能就含有一些敏感信息因此如果EAgent能访问客户端的Skill Library,这可能会导致隐私的泄露。可是如果要对客户端的Skill进行这样的个性化修改好像必须要访问访问客户端的skill吧？我也不太清楚

