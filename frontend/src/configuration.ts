import { Configuration, ConfigurationParameters, DefaultApiFactory } from "./api";

export const apiConfig = new Configuration(
    {
        basePath: "http://localhost:8080/api"
    }
)
export const client = DefaultApiFactory(apiConfig)