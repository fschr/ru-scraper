#!/usr/bin/fish

# Sort all the words into categories (noun, verb, adjective, adverb, numeral, other)
if not test -d data
	mkdir data
end
for i in (seq 2 7)
	switch $i
		case 2
			set filename nouns.txt
		case 3
			set filename verbs.txt
		case 4
			set filename adjectives.txt
		case 5
			set filename adverbs.txt
		case 6
			set filename numerals.txt
		case 7
			set filename other.txt
	end
	cat list.txt | sed '1d' | cut -f$i | grep -v -e "^\$" > data/$filename
end

# cat verbs.txt | sed -n '1!p' | head | xargs -n 1 -i ru-scraper -word \{\} -section Verb
