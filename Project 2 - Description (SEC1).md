Alice, Bob and Charlie are participating in a medical experiment in their local Hospital where they need to provide data for training a Machine Learning (ML) model, which will later be used to diagnose patients. In this setting, the Hospital runs a central server that collects data from all patients.

The Hospital is trusted to execute a data collection protocol established by the Hospital and pre-agreed among all patients taking part in the experiment, without actively trying to cheat on patients and extract and/or reveal their data. However, since the experiment deals with specially sensitive personal data (i.e. health data), neither the patient nor the hospital want plaintext patient data to be directly processed and/or stored by the Hospital. Moreover, the patients do not know each other, and consequently do not trust each other to process and/or store their data. Nevertheless, the patients do trust each other to follow the protocol established by the Hospital. As usual, the patients and the Hospital communicate over insecure networks (i.e. the Internet), where eavesdroppers may try to obtain private data and/or compromise the accuracy of the experiment by tampering with data while it is transmitted.

The restrictions in the scenario above leave the Hospital and patients in a tight spot, since standard ML algorithms require access to plaintext data in order to train a model. Luckily, the researchers in the hospital are collaborating with a team of data scientists and security experts who are all aware of the latest advances in Federated Learning, which allows for distributed training of models using data held locally by different users.
These experts have designed a Federated Learning algorithm that supports a technique called Secure Aggregation, which allows the algorithm to train a model from aggregate data encoding all the plaintext data individually held by each patient into a single value. In particular, in this scenario, it is sufficient for the patient data to be aggregated by summing all individual patient's values into one final aggregate value, which needs to be made available to the Hospital.
Using this technique, neither the patients nor the Hospital get access to each patient's individual plaintext data.  

Design and implement a solution that allows for the 3 patients interacting among each other and with the Hospital in the scenario described above to compute an aggregate value of their private input values in such a way that only the Hospital learns this aggregate value.
Consider that all individual values held by patients are integers in a range [0,...,R] and that the aggregated value is the sum of all individual values, which is also assumed to be in the same range. You must describe an adversarial model (or threat model) capturing the attacks by an adversary who behaves according to the scenario described above, explain how your solution works and why it guarantees security against such an adversary.

Hint: Secure Aggregation for Federated Learning is a real-world practical technique.

Deliverables:  

- A written report describing the adversarial model you are working on, describing the building blocks of your proposed solution, how they are combined in your final solution and why they guarantee security against the adversary you describe.
- An implementation of your solution in a programming language of your choosing, along with clear instructions on how to compile and run it (potentially added to the report or to a separate Readme file).

Submission Instructions:

- Submit only the PDF file with your report and the file containing your code. Do not submit whole folders containing metadata, auxiliary IDE files or anything else than the code and report.
- Please name your submission clearly using your Name/Student ID, e.g. Jane Doe - 36476832.zip, Jane Doe - 36476832 - Report.PDF, Jane Doe - 36476832 - code.c, Jone Doe - 36476832 - Readme.txt. This makes grading faster, so that you get feedback on your hand-in faster.
