import http from 'k6/http';
import { Trend } from 'k6/metrics';
import { group, check, sleep } from "k6";
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

const token = open('/tmp/token');
const BASE_URL = `${__ENV.BASE_URL}`;
const SLEEP_DURATION = 0.2;
const listTrend = new Trend('List_API');
const postTrend = new Trend('POST_API');
const getTrend = new Trend('GET_ID_API');
const patchTrend = new Trend('PATCH_API');

export const options = {
  scenarios: {
    k6_rhtrex: {
      executor: 'constant-arrival-rate',
      rate: 1,
      duration: '300s',
      preAllocatedVUs: 5,
    },
  },
  thresholds: {
    iteration_duration: ['med<900'],
    List_API: ['med<10'],
    POST_API: ['med<10'],
    PATCH_API: ['med<20'],
    GET_ID_API: ['med<10'],
  },  
};


export function handleSummary(data) {
  let formattedStartTime = Math.trunc((new Date().getTime()-30870.873144));
  let formattedEndTime = new Date().getTime();
  data.startTime = formattedStartTime;
  data.endTime = formattedEndTime;
  data["$schema"] = "uri:k6:0.1";
  return {
    '/workspace/output.json': JSON.stringify(data), //the default data object
  };
}

export default function () {
  var id = ''
  const options = {
    headers: {
      "Content-Type": "application/json",
      Authorization: "Bearer " + token,
    },
  };
    group("/api/rh-trex/v1/dinosaurs", () => {
        let search = ''; // specify value as there is no example value for this parameter in OpenAPI spec
        let size = ''; // specify value as there is no example value for this parameter in OpenAPI spec
        let orderBy = ''; // specify value as there is no example value for this parameter in OpenAPI spec
        let page = ''; // specify value as there is no example value for this parameter in OpenAPI spec
        let fields = ''; // specify value as there is no example value for this parameter in OpenAPI spec

        // Request No. 1: 
        {
            let url = BASE_URL + `/api/rh-trex/v1/dinosaurs?page=${page}&size=${size}&search=${search}&orderBy=${orderBy}&fields=${fields}`;
            let request = http.get(url,options);

            check(request, {
                "A JSON array of dinosaur objects": (r) => r.status === 200
            });
            listTrend.add(request.timings.duration);
            sleep(SLEEP_DURATION);

        }
        {
            let url = BASE_URL + `/api/rh-trex/v1/dinosaurs`;
            let body = { "species":  randomString(8)};
            let request = http.post(url, JSON.stringify(body), options);

            id=request.json('id')

            check(request, {
                "Created": (r) => r.status === 201
            });

            postTrend.add(request.timings.duration);
            sleep(SLEEP_DURATION);
        }

    });

      group("/api/rh-trex/v1/dinosaurs/{id}", () => {

        // Request No. 1: 
        {
            let url = BASE_URL + `/api/rh-trex/v1/dinosaurs/${id}`;
            let request = http.get(url,options);

            check(request, {
                "Dinosaur found by id": (r) => r.status === 200
            });

            getTrend.add(request.timings.duration);
            sleep(SLEEP_DURATION);

        }
        {
            let url = BASE_URL + `/api/rh-trex/v1/dinosaurs/${id}`;
            let body = {"species": randomString(8)};
            let request = http.patch(url, JSON.stringify(body), options);

            check(request, {
                "Dinosaur updated successfully": (r) => r.status === 200
            });


            patchTrend.add(request.timings.duration);
            sleep(SLEEP_DURATION);
        }


    });
}
