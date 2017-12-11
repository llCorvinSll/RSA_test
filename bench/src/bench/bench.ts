import Axios from "axios";
import {JSEncrypt} from "jsencrypt";

const HOST_URL = "http://localhost:8080";

export function runBench() {
    singleRun();
}

function singleRun() {
    let PUB_KEY = "";
    let UUID = "";
    let count = 0;

    Axios.get(`${HOST_URL}/test/start`).then((data) => {
        console.log(data.data);


        PUB_KEY = data.data.public_key;
        UUID = data.data.test_id;

        return encriptEndWerify(PUB_KEY, UUID);

    });
}

function encriptEndWerify(pub_key:string, uuid:string):Promise<{}> {
    return Axios.get(HOST_URL + "/test/data/" + uuid).then((data) => {
        console.log(data.data);

        console.log("jsencrypt", JSEncrypt);

        let encryptor = new JSEncrypt();

        console.log(encryptor)

        return {};
    });
}
