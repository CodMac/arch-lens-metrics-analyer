package metrics.loc;

public class MultiLineStrings {
    public String s = """
            This is
            a multi-line
            string
            """;
}
/*
L1: package (1)
L2: (empty)
L3: public class (2)
L4: public String (3)
L5: This is (4)
L6: a multi-line (5)
L7: string (6)
L8: ""; (7)
L9: } (8)
Java Text Blocks: everything inside is considered content.
*/
