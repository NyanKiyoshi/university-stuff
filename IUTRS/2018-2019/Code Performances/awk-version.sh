#!/bin/bash
# search_optimise.sh repertoire mot ...
# recherche (recursivement !) les fichiers du repertoire qui contiennent un des
# mots en argument.
# affiche pour chaque fichier: son chemin absolu et la liste des mots qu'il
# contient

[ $# -ne 1 ] && { echo "usage: $(basename $0) repertoire..." 1>&2; exit 1; }
[ -d "$1" ] || { echo "$(basename $0): $1: repertoire non valide !" 1>&2; exit
1; }

dir=$1
abs_path=$(cd $dir ; pwd)

# Patterns to flag as invalid
patterns="0.0.0.0
127.0.0.1
255.255.255.255
AF_INET
gethostbyaddr
gethostbyname
gethostbyname_ex
Inet4Address
inet_aton
inet_nton
sockaddr_in"

# Join the patterns in a group of whole word matching regexes
# First, it replaces dots with \.
# Then, it groups the patterns (but the <$entry> syntax is only for gawk)
re_patterns=$(echo -e "$patterns" | sed 's/\./\\./g' | paste -sd '|')

declare -A occurrences

# Create the hashmap
for pattern in $patterns; do
    occurrences[$pattern]=""
done

# Retrieve the hashmap keys
KEYS=(${!occurrences[@]})

# usage: PATTERN FILE_PATH
append_occurrence() {
    occurrences[$1]="${occurrences[$1]}$2
"  # <- note the newline
}

# usage: VAL1 VAL2
print_padded() {
    printf "%-20s %s\n" "$1" "$2"
}

print_occurrences() {
    IFS="
"
    for (( I=0; $I < ${#occurrences[@]}; I+=1 )); do
        found_pattern=${KEYS[$I]}
        file_list=${occurrences[$found_pattern]}

        if [ -z "$file_list" ]; then
            print_padded "$found_pattern" Success
        else
            print_padded "$found_pattern" Failed
            for path in $file_list; do
                print_padded + "$path"
            done
        fi
    done
}

find $abs_path -type f \( -name '*.h' -o -name '*.c' -o -name '*.cpp' \) | \
    xargs awk "
BEGIN {
    is_in_comment_body = 0;
    previous_filename = \"\";
}

{
    if (is_in_comment_body) {
        end_index = index(\$0, \"*/\");
        if (end_index) {
            is_in_comment_body = 0;
            \$0 = substr(\$0, end_index + 2, length(\$0));
        }
    }

    if (previous_filename != FILENAME) {
        is_in_comment_body = 0;  # reset
        previous_filename = FILENAME;
    }

    if (!is_in_comment_body) {
        line = \$0;

        // Remove any comments
        gsub(\"\\s*//.*\", \"\", line)

        # Look if there is a comment block comment, which is: /*
        comment_block_index = index(\$0, \"/*\"); 

        # If there is a block comment, remove it 
        # and flag we are in a comment body
        if (comment_block_index) { 
            line = substr(\$0, 0, comment_index - 1);

            after_comment = substr(\$0, comment_index, length(\$0));
            end_comment_block_index = index(after_comment, \"*/\");

            if (!end_comment_block_index) {
                is_in_comment_body = 1;
            }
        }

        if (line) {
            m = match(line, /${re_patterns}/);
            if (m) {
                m = substr(line, RSTART, RLENGTH);
                print sprintf(\"%s:%s\", FILENAME, m);
            }
        }
    }
}" | sort -u | \
    {
        IFS=":"
        while read file_path pattern_found; do
            append_occurrence $pattern_found $file_path
        done
        print_occurrences
    }

