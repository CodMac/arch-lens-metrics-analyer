package com.app.core.module1;

import com.app.core.module1.OtherLocal; // Same domain
import com.app.core.module2.Target2; // Domain 2
import com.app.core.module3.Target3; // Domain 3

public class SourceFile {
}
// Expected CDC: 2 (Target2 belongs to module2, Target3 belongs to module3)
