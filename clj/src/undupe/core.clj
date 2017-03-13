(ns undupe.core
  (:use [clojure.java.io :as io]
        [clojure.pprint]
        [clojure.tools.trace])
  (:require [digest])
  (:gen-class))

(defn walkf [root]
  (filter #(.isFile %) (file-seq (io/file root))))

(defn gather [g [x & xs]]
  "we want to return the set being gathered to. use g as accumulator."
  (when (complement (nil? x))
    (do
      (prn x g ((last x) g))
      (conj g `((first x) ((last x) g))))))


(defn lazy-sieve [s]
  (cons (first s)
    (lazy-seq (lazy-sieve (remove #(zero? (mod % (first s))) (rest s))))))


(defn gen-iterator-keyvals [iterators prime]
  (mapcat #(list (+ % prime) [%]) iterators))

(defn update-iterators [prime iterator-map]
  (let [iterators (apply hash-map (gen-iterator-keyvals (get iterator-map prime) prime))
        basemap (dissoc iterator-map prime)]
    (merge-with concat basemap iterators {(* prime prime) [prime]})))

(defn lazy-erastosthenes [iterator-map [x & xs]]
  (if (contains? iterator-map x)
    (lazy-erastosthenes (update-iterators x iterator-map ) xs)
    (cons x (lazy-seq (lazy-erastosthenes (merge-with concat iterator-map {(* x x) [x]}) xs)))))

(defn primes []
  (cons 2 (lazy-seq (lazy-erastosthenes {4 [2]} (iterate inc 2) ))))

(defn oprimes []
  (lazy-sieve (iterate inc 2) ))

(defn canonical-sieve [s]
  (when (seq s)
    (cons (first s)
          (lazy-seq (canonical-sieve (filter #(not= 0 (mod % (first s)))
                       (rest s))))))
  )

(defn cprimes [] (canonical-sieve (iterate inc 2)))

(take 3 (primes))


(lazy-seq [1 2 3])

(defn sf [n]
  #(zero? (mod % n)))

(filter (comp (sf 2) (sf 5)) (range 15))
(filter #(even? %) (range 15))

()