// Copyright [2022] <liushenggen>

#ifndef CRE2_CGO_H_
#define CRE2_CGO_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "./cre2.h"

bool match(cre2_regexp_t *rex, const char *text, int textlen);

int all_matches(cre2_regexp_t *rex, const char *text, int textlen,
                cre2_string_t *match, int nmatch, int nsubmatch);

int all_matches_index(cre2_regexp_t *rex, const char *text, int textlen,
                      int *match, int nmatch, int nsubmatch);

#ifdef __cplusplus
}  // extern "C"
#endif

#endif  // CRE2_CGO_H_
