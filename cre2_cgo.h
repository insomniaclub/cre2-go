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

    CRE2_CGO_DECL int all_matches(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t *match, int nmatch, int nsubmatch);
    CRE2_CGO_DECL bool match(cre2_regexp_t *rex, const char *text, int textlen);
    CRE2_CGO_DECL int find_all_string_index(cre2_regexp_t *rex, const char *text, int textlen, int **match, int nmatch);
    CRE2_CGO_DECL int find_all_string_submatch_index(cre2_regexp_t *rex, const char *text, int textlen, int **match, int nmatch, int nsubmatch);

#ifdef __cplusplus
} // extern "C"
#endif

#endif