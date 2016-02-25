
#include <stdio.h>
#include <string.h>

#define con_val(arg, var, fmt)                          \
    if(!strcasecmp(_arg, arg)){                         \
        sscanf(_val, fmt, &var);                        \
    }

#define con_match(id)                                   \
    if(!strcasecmp(_val, #id)){                         \
        _v = id;                                        \
    }
        
#define con_enum(arg, var, match_cases...)              \
    if(!strcasecmp(_arg, arg)){                         \
        typeof(var) _v;                                 \
        match_cases                                     \
        var = _v;                                       \
    }
    
    
#define con_load(file, cases...)                        \
{                                                       \
    FILE* _f = fopen(file, "r");                        \
    char _line[64] = {0};                               \
    while(fgets(_line, 64, _f)){                        \
        if(!strncmp(_line, "--", 2)){                   \
            char _arg[32];                              \
            char _val[32];                              \
            sscanf(_line, "--%s %s\n", _arg, _val);     \
            cases                                       \
        }                                               \
    }                                                   \
}

