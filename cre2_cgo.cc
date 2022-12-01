// Copyright [2022] <liushenggen>

#include "./cre2_cgo.h"

#include "./cre2.h"

bool match(cre2_regexp_t *rex, const char *text, int textlen) {
  return cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, NULL, 0) ==
         1;
}

int all_matches(cre2_regexp_t *rex, const char *text, int textlen,
                cre2_string_t *match, int nmatch, int nsubmatch) {
  int cnt = 0;
  while (cnt < nmatch && textlen > 0) {
    cre2_string_t *submatch =
        reinterpret_cast<cre2_string_t *>(match + cnt * nsubmatch);
    if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, submatch,
                    nsubmatch)) {
      return cnt;
    }
    textlen -= submatch->data + submatch->length - text;
    text = submatch->data + submatch->length;
    cnt++;
  }
  return cnt;
}

int all_matches_index(cre2_regexp_t *rex, const char *text, int textlen,
                      int *match, int nmatch, int nsubmatch) {
  int cnt = 0;
  const char *start_addr = text;
  while (cnt < nmatch && textlen > 0) {
    cre2_string_t str[nsubmatch];
    if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, &str[0],
                    nsubmatch)) {
      return cnt;
    }
    for (int i = 0; i < nsubmatch; i++) {
      reinterpret_cast<int *>(match + cnt * nsubmatch * 2)[2 * i] =
          str[i].data - start_addr;
      reinterpret_cast<int *>(match + cnt * nsubmatch * 2)[2 * i + 1] =
          ((int *)match + cnt * nsubmatch * 2)[2 * i] + str[i].length;
    }
    textlen -= str->data + str->length - text;
    text = str->data + str->length;
    cnt++;
  }
  return cnt;
}
