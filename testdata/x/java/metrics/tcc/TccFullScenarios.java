package com.test.tcc;

/**
 * 场景 1: 高内聚 (High Cohesion)
 * NP = 3 * 2 / 2 = 3 (Pairs: M1-M2, M1-M3, M2-M3)
 * Shared Fields:
 * - (M1, M2) share f1
 * - (M2, M3) share f2
 * - (M1, M3) share f1 (indirect) or f2 (indirect)? 
 *   Wait, TCC definition: "directly connected if they access at least one common instance variable".
 *   - M1(f1), M2(f1, f2), M3(f2)
 *   - Pairs: (M1,M2) via f1, (M2,M3) via f2, (M1,M3) via none.
 *   - NDP = 2. TCC = 2/3 = 0.66
 */
class HighCohesion {
    private int f1, f2;
    public void m1() { f1 = 1; }
    public void m2() { f1 = 2; f2 = 2; }
    public void m3() { f2 = 3; }
}

/**
 * 场景 2: 黄金判定点 (TCC = 0.33)
 * NP = 3 (M1, M2, M3)
 * NDP = 1 (Only M1 and M2 share a field)
 * TCC = 1/3 = 0.33
 */
class MediumCohesion {
    private int f1, f2, f3;
    public void m1() { f1 = 1; }
    public void m2() { f1 = 2; }
    public void m3() { f3 = 3; } // Isolated
}

/**
 * 场景 3: 零内聚 (Zero Cohesion)
 * 每个方法访问独立的字段
 * NDP = 0, TCC = 0
 */
class LowCohesion {
    private int f1, f2, f3;
    public void m1() { f1 = 1; }
    public void m2() { f2 = 2; }
    public void m3() { f3 = 3; }
}

/**
 * 场景 4: 单方法类 (Edge Case)
 * 根据实现，TCC 应返回 1.0
 */
class SingleMethod {
    private int f1;
    public void m1() { f1 = 1; }
}
