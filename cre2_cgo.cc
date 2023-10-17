#include <stdlib.h>

#include "cre2.h"
#include "cre2_cgo.h"

int all_matches(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t *match, int nmatch, int nsubmatch)
{
    int cnt = 0;
    while (cnt < nmatch && textlen > 0)
    {
        cre2_string_t *submatch = (cre2_string_t *)match + cnt * nsubmatch;
        if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, submatch, nsubmatch))
        {
            return cnt;
        }
        textlen -= submatch->data + submatch->length - text;
        text = submatch->data + submatch->length;
        cnt++;
    }
    return cnt;
}

bool match(cre2_regexp_t *rex, const char *text, int textlen)
{
    return cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, NULL, 0) == 1;
}

int find_all_string_index(cre2_regexp_t *rex, const char *text, int textlen, int **match, int nmatch)
{
    int cnt = 0;
    const char *start_addr = text;
    while (cnt < nmatch && textlen > 0)
    {
        cre2_string_t str;
        if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, &str, 1))
        {
            return cnt;
        }
        ((int *)match + cnt * 2)[0] = str.data - start_addr;
        ((int *)match + cnt * 2)[1] = ((int *)match + cnt * 2)[0] + str.length;
        textlen -= str.data + str.length - text;
        text = str.data + str.length;
        cnt++;
    }
    return cnt;
}

int find_all_string_submatch_index(cre2_regexp_t *rex, const char *text, int textlen, int **match, int nmatch, int nsubmatch)
{
    int cnt = 0;
    const char *start_addr = text;
    cre2_string_t *strs = (cre2_string_t *)malloc(sizeof(cre2_string_t)*nsubmatch);
    while (cnt < nmatch && textlen > 0)
    {
        if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, strs, nsubmatch))
        {
            break;
        }
        for (int i = 0; i < nsubmatch; i++) {
            if ((strs+i)->data == NULL) {
                ((int *)match + cnt * 2*nsubmatch)[i*2] = 0;
                ((int *)match + cnt * 2*nsubmatch)[i*2+1] = 0;
                continue;
            }
            ((int *)match + cnt * 2*nsubmatch)[i*2] = (strs+i)->data - start_addr;
            ((int *)match + cnt * 2*nsubmatch)[i*2+1] = ((int *)match + cnt * 2*nsubmatch)[i*2] + (strs+i)->length;
        }
        textlen -= strs->data + strs->length - text;
        text = strs->data + strs->length;
        cnt++;
    }
    free(strs);
    strs = NULL;
    return cnt;
}