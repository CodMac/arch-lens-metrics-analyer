package com.test.wmc;

/**
 * 场景 1: 简单类
 * WMC = 1 (m1) + 1 (m2) + 1 (default constructor) = 3
 */
class SimpleClass {
    public void m1() {}
    public void m2() {}
}

/**
 * 场景 2: 上帝类候选者 (WMC > 47)
 * 通过高复杂度的逻辑节点累加
 */
class GodClassCandidate {
    
    // 基础复杂度 1
    // + 5 (if) + 5 (for) = 11
    public void heavyMethod(int x) {
        if (x > 0) {}
        if (x > 1) {}
        if (x > 2) {}
        if (x > 3) {}
        if (x > 4) {}

        for (int i=0; i<10; i++) {}
        for (int i=0; i<10; i++) {}
        for (int i=0; i<10; i++) {}
        for (int i=0; i<10; i++) {}
        for (int i=0; i<10; i++) {}
    }

    // 静态块复杂度贡献
    // 1 (base) + 10 (if) = 11
    static {
        int x = 10;
        if (x > 0) {} if (x > 1) {} if (x > 2) {} if (x > 3) {} if (x > 4) {}
        if (x > 5) {} if (x > 6) {} if (x > 7) {} if (x > 8) {} if (x > 9) {}
    }

    // 另一个复杂方法
    // 1 (base) + 25 (case) = 26
    public void switchMethod(int x) {
        switch(x) {
            case 1:case 2:case 3:case 4:case 5:
            case 6:case 7:case 8:case 9:case 10:
            case 11:case 12:case 13:case 14:case 15:
            case 16:case 17:case 18:case 19:case 20:
            case 21:case 22:case 23:case 24:case 25:
                break;
        }
    }

    // Total WMC = 11 + 11 + 26 + 1 (implicit constructor) = 49 (Trigger > 47)
}
