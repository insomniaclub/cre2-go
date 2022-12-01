#include <string.h>
#include <stdio.h>
#include <stdarg.h>

#include "cre2.h"
#include "cre2_cgo.h"

int all_matches(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t *match, int nmatch, int nsubmatch)
{
    for (int i = 0; i < nmatch && textlen > 0; i++)
    {
        cre2_string_t *submatch = (cre2_string_t *)match + i * nsubmatch;
        if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, submatch, nsubmatch))
        {
            return i;
        }
        text = submatch->data + submatch->length;
        textlen -= submatch->length;
    }
    return nmatch;
}

bool match_string(cre2_regexp_t *rex, const char *text, int textlen)
{
    return cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, NULL, 0) == 1;
}

int find_all_string(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t *match, int nmatch)
{
    return all_matches(rex, text, textlen, match, nmatch, 1);
}

int find_all_string_submatch(cre2_regexp_t *rex, const char *text, int textlen, cre2_string_t **match, int nmatch)
{
    return all_matches(rex, text, textlen, (cre2_string_t *)match, nmatch, cre2_num_capturing_groups(rex) + 1);
}

int find_all_string_index(cre2_regexp_t *rex, const char *text, int textlen, int **match, int nmatch)
{
    const char *start_addr = text;
    for (int i = 0; i < nmatch && textlen > 0; i++)
    {
        cre2_string_t str;
        if (!cre2_match(rex, text, textlen, 0, textlen, CRE2_UNANCHORED, &str, 1))
        {
            return i;
        }
        ((int *)match + i * 2)[0] = str.data - start_addr;
        ((int *)match + i * 2)[1] = ((int *)match + i * 2)[0] + str.length;
        text = str.data + str.length;
        textlen -= str.length;
    }
    return nmatch;
}
