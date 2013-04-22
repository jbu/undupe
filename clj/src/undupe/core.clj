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
      (conj g `((first x) ((last x) g))))
  ))

(defn ^:dynamic sieve [primes xs]
  (if-let [prime (first xs)]
    (recur (conj primes prime) (remove #(zero? (mod % prime)) xs))
    primes))

(def ^:dynamic rsieve
  (cons n (lsieve (inc n))))

(defn -main
  [& args]
  (dotrace [sieve] (sieve [1] (range 2 10)))
  ; (pprint  (pmap #(`(% (keyword (digest/md5 %)))) (mapcat #(walkf %) args)))
  ; (shutdown-agents)
    )
