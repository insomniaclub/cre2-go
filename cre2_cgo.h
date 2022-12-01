#ifndef CRE2_CGO_H
#define CRE2_CGO_H

#ifdef __cplusplus
extern "C"
{
#endif

#include "cre2.h"

#ifndef CRE2_CGO_DECL
#define CRE2_CGO_DECL extern
#endif

    //  match_string return if the string matches the regular expression
    CRE2_CGO_DECL bool match_string(cre2_regexp_t *rex, const char *text, int textlen);
    //  find_all_substring match all the string which satisfied the regular expression and return the number of matches
    CRE2_CGO_DECL int find_all_string(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t *match, int nmatch);
    //  find_all_substring_index match all the begin of string which satisfied the regular expression and return the number of matches
    CRE2_CGO_DECL int find_all_string_index(cre2_regexp_t *rex, const char *text, int textlen, int **match, int nmatch);
    CRE2_CGO_DECL int find_all_string_submatch(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t **match, int nmatch);

#ifdef __cplusplus
} // extern "C"
#endif

#endif