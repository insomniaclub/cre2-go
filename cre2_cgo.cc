#include "cre2.h"
#include "cre2_cgo.h"

bool match(cre2_regexp_t *rex, const char *text, int textlen)
{
    return cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, NULL, 0) == 1;
}

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
        text = str.data + str.length;
        textlen -= str.length;
        cnt++;
    }
    return cnt;
}
