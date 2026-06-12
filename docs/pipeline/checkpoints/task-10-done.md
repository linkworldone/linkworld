# Task 10 — T9 Billing/BillDetail（付账两步态+手续费读链+paying+修双重缩放）（web 3/3）

> 子项目 web(3/3) | 状态 DONE | 2026-06-10

### 产出文件
- packages/web/src/hooks/contracts/useFeeManager.ts（新建：useFeeRate 读 getFeeRate 基点/10000 + useCalculateFee 读 calculateFee 不自算；失败 isError 禁写死兜底）
- packages/web/src/components/shared/FeeBreakdown.tsx（新建：手续费明细 loading→skeleton/失败→--）
- packages/web/src/pages/{Billing,BillDetail}.tsx（重写：修双重缩放、TwoStepAction 付账、paying 态、FeeBreakdown、lucide+ui/card+AmountDisplay）
- packages/web/src/hooks/useBilling.ts（recordIntent 拒签上抛对齐 useDeposit）
- packages/web/src/hooks/contracts/usePaymentContract.ts（payBill 去 _value 第二参）+ index.ts 导出 fee hooks
- packages/web/src/config/constants.ts（删 PLATFORM_FEE_RATE）
- 3 个新测试

### git commit
7b2a7f0 feat: web T9 Billing/BillDetail(付账两步态+手续费读链+paying+修 totalAmount 双重缩放+深蓝金)

### TDD
先红后绿：FEE/BILL 测先写 → 实现后 91 passed。

### 测试结果
npx tsc --noEmit 0 error。npm run build ✓。npx vitest run：20 files / 91 passed（前序 78+新 13，无回归）。用例 FEE-01(getFeeRate 基点→1.5%/calculateFee 不自算)/FEE-02(读链失败→-- 不写死)/BILL-02(展示 formatAmount(6)=5.00 DOM 不含 5000000)/BILL-03(paying TxStatusBadge 不染绿不据 200 置已付)。grep 无 formatUSD/parseUnits/2.5%/emoji。主 Agent 已独立复跑确认 tsc 0+build ✓+91 测。

### code-simplifier
useFeeRate/useCalculateFee 读链复用；FeeBreakdown 组件供 T10 复用；删 PLATFORM_FEE_RATE 死常量。

### spec review
按 design v2 §3.2/§3.3/§3.6 + arch-review B2/对账不据 200 + T2 资损遗留执行。**核心资损修复**：totalAmount 6 位最小单位不再 parseUnits 二次缩放（×10^6 错百万倍）；付账额最小单位 bigint、授权额=本金+calculateFee 读链。paying TxStatusBadge 禁绿、is_paid 等后端事件回填。未越界（RegionDetail T10/其余 T11、接链基建只调用）。

### 设计还原
Billing/BillDetail 对齐 design §3.2 付账两步/§3.6 手续费读链/§3.3 paying；深蓝金 ui/card+lucide+AmountDisplay。

### 复用检查
复用 TwoStepAction(T3)/TxStatusBadge·AmountDisplay(T6)/signedPost(T5)/FeeManager ABI calculateFee/formatAmount(6)；新增 useFeeRate/useCalculateFee/FeeBreakdown。

### 设计稿对照
数值对照：手续费率读链 getFeeRate 150/10000=1.5%（FEE-01）✅；付账额=本金+calculateFee 读链不自算 ✅；totalAmount 展示 5.00 不二次缩放（BILL-02）✅；paying 不染绿（BILL-03）✅；读链失败 --（FEE-02）✅；91 测/tsc 0/build ✓ ✅。

### 新增组件
新增 useFeeRate/useCalculateFee/FeeBreakdown。

### 新增色值
无（用 T6 深蓝金语义类）。

### ⚠️ 遗留（带入 T10/T11）
- T10 RegionDetail 申请号码弹层手续费可直接复用 FeeBreakdown+useFeeRate/useCalculateFee（design §3.6 同源）。
- useBills 拉全量客户端按 Tab 过滤；useBilling 仍导出 isContractPending/isConfirming/isSuccess（页面已不用，T11 可清）。
- format.ts 的 formatUSD 现无引用（dead util），删除可选 T11 清理。
- operatorApi.requiredDeposit parseEther(18 位) T2 遗漏仍在 → T10 RegionDetail 换肤时同步修 6 位。
