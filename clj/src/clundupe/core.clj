(ns clundupe.core
  (:use [clojure.java.io :as io]
        [clojure.pprint])
  (:require [digest])
  (:gen-class))

(defn walkf [root]
  (filter #(.isFile %) (file-seq (io/file root))))

(defn -main
  "I don't do a whole lot."
  [& args]
  (pprint (pmap #(vector % (keyword (digest/md5 %))) (mapcat #(walkf %) args)))
  (shutdown-agents)
    )

