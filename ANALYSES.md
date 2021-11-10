# Some Latin Texts Analysed and Compared with riglib.

# Terms.
The **riglib score (RS)** is defined as the unique word-count of the riglib output of a document,
and the **riglib overlap (RO)** of two documents is the word-count of the intersection of their
outputs. For more on this, see the [readme](README.md).

# Establishing benchmarks.
The RS and RO (given two texts) figures are difficult to apply directly to the learning process, so
several benchmark-texts shall be selected. 

As a specimen of the words that intermediate learners might know, I have selected the 
[world list](https://www.hackettpublishing.com/pdfs/Familia_Romana_Latin-English_Vocabulary.pdf)
from _Familia Romana_ of _Lingua Latina pe se Illustrata_
and also the [DCC Latin core vocabulary](https://dcc.dickinson.edu/latin-core-list1) list.

For theology I have chosen Anselm's _Proslogion_ because (as an intermediate learner) I found
it so accessible that part of my initial motivation for building riglib was the desire to assess the
difficulty of other works relative to it. To this I have added _De Gratia et Libero Arbitrio_
and _Retractiones_ (Augustine), and _De Carne Christi_ (Tertullian).

| Work                          | RS   |
| ----------------------------- | ---- |
| Familia Romana                | 2323 |
| DCC Core List                 | 1814 |
| **Theology**                  |      |
| Proslogion                    | 1399 |
| De Gratia et Libero Arbitrio  | 2020 |
| Retractiones                  | 3781 |
| Confessiones                  | 8492 |
| **Classical**                 |      |
