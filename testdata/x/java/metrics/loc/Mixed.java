package metrics.loc;

public class Mixed {
    public void test() {
        int a = 1; // code and comment
        /* comment */ int b = 2;
        int c = 3; /* multi-line
                      comment on same line */
    }
}
/*
L1: package (1)
L2: (empty)
L3: public class (2)
L4: public void (3)
L5: int a (4)
L6: int b (5)
L7: int c (6)
L8: (comment part) - L8 is "                      comment on same line \*\/" - this is PURE comment line.
L9: } (7)
L10: } (8)
Let me re-check L7/L8.
L7: int c = 3; /* multi-line
L8:               comment on same line \*\/
L7 has code. L8 has NO code.
Logical LOC: 8
*/
