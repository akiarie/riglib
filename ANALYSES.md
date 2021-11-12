# Some Latin Texts Analysed and Compared with riglib.

# Terms.
The **riglib score (RS)** is defined as the unique word-count of the riglib output of a document,
and the **riglib overlap (RO)** of two documents is the word-count of the intersection of their
outputs. The **known-word density (KWD)**, is defined on two documents, calling the first the
_dictionary_ and the second the _book_, as the ratio of _Whitaker-known_ word-forms in the book
which are present in the dictionary.
For more on these definitions, see the [readme](README.md).

# Choosing benchmarks.
The RS, RO and KWD figures are difficult to apply directly to the learning process, so several
benchmark-texts shall be selected. 

As a specimen of the words that intermediate learners might know, I have selected the 
[world list](https://www.hackettpublishing.com/pdfs/Familia_Romana_Latin-English_Vocabulary.pdf)
from _Familia Romana_ of _Lingua Latina pe se Illustrata_
and also the [DCC Latin core vocabulary](https://dcc.dickinson.edu/latin-core-list1) list.

The works chosen are broken into the categories of _theology_ and _classics_. Also, many of the
works are listed multiple times because I have broken them down into various sections to get insight
into the difficulties of moving within the work, so that a reader studying Augustine, for example,
may know whether studying the _Confessiones_ in sequence is a practical programme.

For theology, aside from the Vulgate, I have chosen Anselm's _Proslogion_ because (as an
intermediate learner) I found it so accessible that part of my initial motivation for building
riglib was the desire to assess the difficulty of other works relative to it. To this I have added
_Retractiones_ and _Confessiones_ (Augustine), _De Carne Christi_ (Tertullian), and _De Imitatio
Christi_ (à Kempis).

Finally, the classics. We have, as samples of prose, _De Natura Deorum_ (Cicero),
_Ab Urbe Condita_ (Livy), _Bellum Iugurthinum_ (Sallust) and _Annales_ (Tacitus). Then, of
course, we have _Aeneid_ (Virgil), 


| Author     | Work                          | RS    |
|------------|-------------------------------| ----  |
| Francese   | DCC Core List                 | 1814  |
| Ørberg     | Familia Romana                | 2323  |
| **Theology**                               |       |
| Anslem     | Proslogion                    | 1399  |
| Tertullian | De Carne Christi              | 2309  |
| Augustine  | Confessiones Liber I          | 3126  |
| Augustine  | Retractiones                  | 3781  |
| à Kempis   | De Imitatio Christi           | 5052  |
| Jerome     | Epistulae (Vulgata)           | 5712  |
| Jerome     | Novus Testamentum (Vulgata)   | 6812  |
| Augustine  | Confessiones                  | 8492  |
| Jerome     | Vetus Testamentum (Vulgata)   | 9998  |
| Jerome     | Biblia Sacra Vulgata          | 12222 |
| **Classical**                              |       |
| Caesar     | De Bello Gallico Liber I      | 2172  |
| Cicero     | De Natura Deorum Liber I      | 2846  |
| Caesar     | De Bello Gallico              | 5128  |
| Cicero     | De Natura Deorum              | 5554  |
