(ns clundupe.core
  (:use [clojure.java.io :as io]
        [clojure.pprint])
  (:gen-class))

;((defmulti take-files ()) take-files [fseq])

(defn -main
  "I don't do a whole lot."
  [& args]
  (for [f args
    :let [fi (file-seq (io/file f))]]
    (pprint fi)
    ))

