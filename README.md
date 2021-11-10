# riglib.
_A tool for analysing words in a text to assess the difficulty of (Latin) books._

# Overview.

_riglib_ is short for _rigor librī_, and therefore functions as a tool for assessing the difficulty
of a Latin work, to facilitate the learner in charting out a route of works to read in sequence.

It works by feeding in each input line into [Whitaker's
words](http://archives.nd.edu/whitaker/words.htm) and then selecting all the word-forms in the
inflection lines it outputs. 

To get an idea of the kind of analysis that it can produce, see [this document](ANALYSES.md).

# Preparing input for riglib.

riglib accepts as its input a list of words, such as can be produced with a command like
```bash
$ cat [TEXT FILE] | tr ' ' '\n' | tr -cd '[:alpha:]\n' | sort | uniq > word-list.txt
```
where `[TEXT FILE]` is the path to a book.

_Note_. The input book must not contain vowel signs such as _ā_, _ē_, _ī_, _ō_, _ū_. To strip these
add
```bash
$ iconv -f utf8 -t ascii//TRANSLIT
```
to the pipeline above, prior to the command to strip non-alpha characters.

# Usage.

To use the tool, simply type
```bash
$ riglib [FILE 1] [FILE 2] [FILE 3] ...
```
where `[FILE 1]`, `[FILE 2]`, `[FILE 3]`, etc. are the files to be analysed.

riglib ignores unknown words, thereby elminating anything which Whitaker's words does not consider
to be Latin. This is advantageous when working with large files and there is no time to filter out
non-Latin words, but may be a problem for well-transcribed documents. There is therefore the
optional _print-unknown mode_,
```bash
$ riglib -u [FILE]
```
or
```bash
$ riglib --print-unknown [FILE]
```
for situations where one desires to retain all the input words.

# Interpreting riglib's output.

For the sake of efficiency, riglib does not bother to sort its output, nor does it avoid repeat
outputs. Thus to obtain a word-count one will regularly use
```bash
$ riglib [FILE] | sort | uniq | wc -l
```
or something of the like.


## Obtaining a riglib score.

This figure, which we may name the **riglib score**, gives one a very accurate estimate regarding
the number of Latin forms or word-families that the inputted work contains.


## Comparing two texts

To compare the difficulty of two texts, first pass them through riglib and store their outputs. Then
run the following (where `[OUTPUT 1]` and `[OUTPUT 2]` as the paths to the two riglib outputs)
```bash
$ sort [OUTPUT 1] [OUTPUT 2] | uniq -d | wc -l
```
to get a count of the intersection set of forms, the **riglib overlap**.
