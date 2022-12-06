// Copyright [2022] <liushenggen>

#include "./cre2_cgo.h"

bool match(cre2_regexp_t *rex, const char *text, int textlen) {
  return cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, NULL, 0);
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
  const char *pstr = text;
  cre2_string_t submatch[nsubmatch];
  while (cnt < nmatch && textlen > 0) {
    if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED,
                    &submatch[0], nsubmatch)) {
      return cnt;
    }
    for (int i = 0; i < nsubmatch; i++) {
      (match + cnt * nsubmatch * 2)[2 * i] = submatch[i].data - pstr;
      (match + cnt * nsubmatch * 2)[2 * i + 1] =
          submatch[i].data - pstr + submatch[i].length;
    }
    textlen -= submatch->data + submatch->length - text;
    text = submatch->data + submatch->length;
    cnt++;
  }
  return cnt;
}
