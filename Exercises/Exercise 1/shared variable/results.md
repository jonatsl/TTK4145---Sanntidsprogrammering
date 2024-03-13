we have concurrency in the threads, they run simultaneously.
they share the same i variable, so they both increment and decrement i, but when they run at the same time the result does not become 0. 

bruker mutex siden den sikrer at maks en thread kan aksessere en shared resource (i dette tilfellet i).


