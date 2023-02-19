import { useContext } from "react";
import { TranslateContext } from "./TranslateProvider";

export const useTranslate = () => useContext(TranslateContext);
