import * as es6 from "es6-promise";
import {runBench} from "./bench/bench";

es6.polyfill();


runBench();
