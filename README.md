
<img align="right" width="50%" src="./go_to_the_goal.jpg">

## Задание: Разработать API получения списка транзакций на блокчейне Ethereum

1. Разработать АПИ которое будет отправлять список транзакций по полученным фильтрам (ID транзакции, Адрес отправителя, Адрес получателя, Номер блока, Дата отправки транзакции), 
   также должна быть предусмотрена валидация данных полученных в запросе и пагинация, данные должны браться из MongoDB.
```
   Список полей:
    - ID транзакции
    - Адрес отправителя
    - Адрес получателя
    - Номер блока
    - Количество подтверждений 
    - Дата отправки транзакции
    - Отправленная сумма
    - Комиссия
```
2. Дополнительно при запуске приложения должна запускаться функция в go рутине, которая циклично (нужно учитывать лимит 
   запросов etherscan и добавить sleep после каждого блока)будет делать запрос на АПИ etherscan и получать текущий блок 
   и транзакции с него (методы которыми можно получить описаны ниже). После получения транзакций они должны сохранятся 
   в MongoDb, откуда их уже будет получать пользователь через АПИ, а также при получении нового блока должно обновляться
   количество подтверждений на транзакциях которые были захвачены до этого.
   
3. При запуске приложения должны инициализироваться последние 1000 блоков с блокчейн.
4. Данные по транзакциям можно получать из https://docs.etherscan.io/api-endpoints/geth-parity-proxy (будет достаточно 
   методов eth_getBlockByNumber и eth_getTransactionByHash)


#### Требования к выполнению:
```
    Код разместить на github
    Рабочий проект развернуть на heroku (или в любой другой песочнице).
```
### Решения:

Не совсем понятно было где брать данные, а именно:
```
    - Количество подтверждений
    - Дата отправки транзакции
    - Комиссия
```

Потому что в ответе на запрос `eth_getTransactionByHash`:

```json
{
   "jsonrpc":"2.0",
   "id":1,
   "result":{
      "blockHash":"0xf850331061196b8f2b67e1f43aaa9e69504c059d3d3fb9547b04f9ed4d141ab7",
      "blockNumber":"0xcf2420",
      "from":"0x00192fb10df37c9fb26829eb2cc623cd1bf599e8",
      "gas":"0x5208",
      "gasPrice":"0x19f017ef49",
      "maxFeePerGas":"0x1f6ea08600",
      "maxPriorityFeePerGas":"0x3b9aca00",
      "hash":"0xbc78ab8a9e9a0bca7d0321a27b2c03addeae08ba81ea98b03cd3dd237eabed44",
      "input":"0x",
      "nonce":"0x33b79d",
      "to":"0xc67f4e626ee4d3f272c2fb31bad60761ab55ed9f",
      "transactionIndex":"0x5b",
      "value":"0x19755d4ce12c00",
      "type":"0x2",
      "accessList":[
         
      ],
      "chainId":"0x1",
      "v":"0x0",
      "r":"0xa681faea68ff81d191169010888bbbe90ec3eb903e31b0572cd34f13dae281b9",
      "s":"0x3f59b0fa5ce6cf38aff2cfeb68e7a503ceda2a72b4442c7e2844d63544383e3"
   }
}
```

Не нашел этих значений. Поэтому я всю логику работы API приложения построил на основе 
получения блоков по запросу `eth_getBlockByNumber` и в mongodb сохранял все транзакции в этом 
блоке с такими полями:
```
    Список полей:
    
    - ID транзакции
    - Адрес отправителя
    - Адрес получателя
    - Номер блока
    - Отправленная сумма
```

Соответственно, по запросу на api с параметром номер блока з mongodb возвращается список всех транзакций
в этом блоке.

### API:

`GET: https://etherscan-test.herokuapp.com/api?block=`

- `block` - номер блока з Etherscan  


### Пример работы приложения

```
❯ git push heroku main
Enumerating objects: 7, done.
Counting objects: 100% (7/7), done.
Delta compression using up to 8 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (4/4), 328 bytes | 328.00 KiB/s, done.
Total 4 (delta 2), reused 0 (delta 0), pack-reused 0
remote: Compressing source files... done.
remote: Building source:
remote:
remote: -----> Building on the Heroku-20 stack
remote: -----> Using buildpack: heroku/go
remote: -----> Go app detected
remote: -----> Fetching stdlib.sh.v8... done
remote: ----->
remote:        Detected go modules via go.mod
remote: ----->
remote:        Detected Module Name: apietherscan
remote: ----->
remote: -----> Using go1.17.3
remote: -----> Determining packages to install
remote:        
remote:        Detected the following main packages to install:
remote:                 apietherscan/cmd
remote:        
remote: -----> Running: go install -v -tags heroku apietherscan/cmd
remote: apietherscan/api
remote: apietherscan/internal/model
remote: apietherscan/configs
remote: apietherscan/pkg/logger
remote: apietherscan/internal/store
remote: apietherscan/pkg/db
remote: apietherscan/internal/etherscan
remote: apietherscan/internal/handlers
remote: apietherscan/cmd
remote:        
remote:        Installed the following binaries:
remote:                 ./bin/cmd
remote:        
remote:        Created a Procfile with the following entries:
remote:                 web: bin/cmd
remote:        
remote:        If these entries look incomplete or incorrect please create a Procfile with the required entries.
remote:        See https://devcenter.heroku.com/articles/procfile for more details about Procfiles
remote:        
remote: -----> Discovering process types
remote:        Procfile declares types -> web
remote:
remote: -----> Compressing...
remote:        Done: 4.5M
remote: -----> Launching...
remote:        Released v14
remote:        https://etherscan-test.herokuapp.com/ deployed to Heroku
remote:
remote: Verifying deploy... done.
To https://git.heroku.com/etherscan-test.git
411350a..b6e69af  main -> main
```

```
2022-01-30T17:55:01.000000+00:00 app[api]: Build started by user ivnxerinov@gmail.com
2022-01-30T17:55:10.349680+00:00 app[api]: Deploy b6e69af3 by user ivnxerinov@gmail.com
2022-01-30T17:55:10.349680+00:00 app[api]: Release v14 created by user ivnxerinov@gmail.com
2022-01-30T17:55:11.728395+00:00 heroku[web.1]: State changed from down to starting
2022-01-30T17:55:12.081512+00:00 heroku[web.1]: Starting process with command `bin/cmd`
2022-01-30T17:55:13.050530+00:00 app[web.1]: time="2022-01-30T17:55:13Z" level=info msg="Start server API on Port - 52153"
2022-01-30T17:55:13.710079+00:00 heroku[web.1]: State changed from starting to up
2022-01-30T17:55:28.000000+00:00 app[api]: Build succeeded
2022-01-30T18:05:19.180002+00:00 app[web.1]: time="2022-01-30T18:05:19Z" level=info msg="Successful last thousand load block"
2022-01-30T18:05:19.315138+00:00 app[web.1]: time="2022-01-30T18:05:19Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:19.333059+00:00 app[web.1]: Diff block -  14108660
2022-01-30T18:05:20.451728+00:00 app[web.1]: time="2022-01-30T18:05:20Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:20.653937+00:00 app[web.1]: Diff block -  14108662
2022-01-30T18:05:19.912852+00:00 app[web.1]: Diff block -  14108661
2022-01-30T18:05:22.019669+00:00 app[web.1]: Diff block -  14108665
2022-01-30T18:05:21.124575+00:00 app[web.1]: Diff block -  14108663
2022-01-30T18:05:22.605066+00:00 app[web.1]: Diff block -  14108666
2022-01-30T18:05:22.721971+00:00 app[web.1]: time="2022-01-30T18:05:22Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:21.571106+00:00 app[web.1]: Diff block -  14108664
2022-01-30T18:05:21.586991+00:00 app[web.1]: time="2022-01-30T18:05:21Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:23.188765+00:00 app[web.1]: Diff block -  14108667
2022-01-30T18:05:24.212587+00:00 app[web.1]: Diff block -  14108669
2022-01-30T18:05:24.949673+00:00 app[web.1]: Diff block -  14108670
2022-01-30T18:05:24.991385+00:00 app[web.1]: time="2022-01-30T18:05:24Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:25.546783+00:00 app[web.1]: Diff block -  14108671
2022-01-30T18:05:23.771719+00:00 app[web.1]: Diff block -  14108668
2022-01-30T18:05:23.856414+00:00 app[web.1]: time="2022-01-30T18:05:23Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:26.624250+00:00 app[web.1]: Diff block -  14108673
2022-01-30T18:05:27.085141+00:00 app[web.1]: Diff block -  14108674
2022-01-30T18:05:27.381995+00:00 app[web.1]: time="2022-01-30T18:05:27Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:28.273980+00:00 app[web.1]: Diff block -  14108676
2022-01-30T18:05:26.003969+00:00 app[web.1]: Diff block -  14108672
2022-01-30T18:05:26.126592+00:00 app[web.1]: time="2022-01-30T18:05:26Z" level=info msg="Current block number - 14108701"
2022-01-30T18:05:29.219358+00:00 app[web.1]: Diff block -  14108678
2022-01-30T18:05:28.745006+00:00 app[web.1]: Diff block -  14108677
2022-01-30T18:05:28.905913+00:00 app[web.1]: time="2022-01-30T18:05:28Z" level=info msg="Successful load block number - 14108701"
2022-01-30T18:05:28.905918+00:00 app[web.1]: time="2022-01-30T18:05:28Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:30.044633+00:00 app[web.1]: time="2022-01-30T18:05:30Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:30.244459+00:00 app[web.1]: Diff block -  14108680
2022-01-30T18:05:30.695582+00:00 app[web.1]: Diff block -  14108681
2022-01-30T18:05:31.180387+00:00 app[web.1]: time="2022-01-30T18:05:31Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:31.417342+00:00 app[web.1]: Diff block -  14108682
2022-01-30T18:05:32.317200+00:00 app[web.1]: time="2022-01-30T18:05:32Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:27.688722+00:00 app[web.1]: Diff block -  14108675
2022-01-30T18:05:29.664693+00:00 app[web.1]: Diff block -  14108679
2022-01-30T18:05:32.577405+00:00 app[web.1]: Diff block -  14108684
2022-01-30T18:05:34.044337+00:00 app[web.1]: Diff block -  14108687
2022-01-30T18:05:33.162437+00:00 app[web.1]: Diff block -  14108685
2022-01-30T18:05:31.989596+00:00 app[web.1]: Diff block -  14108683
2022-01-30T18:05:33.452141+00:00 app[web.1]: time="2022-01-30T18:05:33Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:33.603251+00:00 app[web.1]: Diff block -  14108686
2022-01-30T18:05:35.219853+00:00 app[web.1]: Diff block -  14108689
2022-01-30T18:05:34.586666+00:00 app[web.1]: time="2022-01-30T18:05:34Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:34.631356+00:00 app[web.1]: Diff block -  14108688
2022-01-30T18:05:36.113718+00:00 app[web.1]: Diff block -  14108691
2022-01-30T18:05:36.691522+00:00 app[web.1]: Diff block -  14108692
2022-01-30T18:05:36.856230+00:00 app[web.1]: time="2022-01-30T18:05:36Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:35.672514+00:00 app[web.1]: Diff block -  14108690
2022-01-30T18:05:35.721475+00:00 app[web.1]: time="2022-01-30T18:05:35Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:37.135744+00:00 app[web.1]: Diff block -  14108693
2022-01-30T18:05:37.990110+00:00 app[web.1]: time="2022-01-30T18:05:37Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:38.015163+00:00 app[web.1]: Diff block -  14108695
2022-01-30T18:05:37.576835+00:00 app[web.1]: Diff block -  14108694
2022-01-30T18:05:38.588661+00:00 app[web.1]: Diff block -  14108696
2022-01-30T18:05:39.600605+00:00 app[web.1]: Diff block -  14108698
2022-01-30T18:05:40.268403+00:00 app[web.1]: time="2022-01-30T18:05:40Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:40.493844+00:00 app[web.1]: Diff block -  14108699
2022-01-30T18:05:39.027734+00:00 app[web.1]: Diff block -  14108697
2022-01-30T18:05:39.127625+00:00 app[web.1]: time="2022-01-30T18:05:39Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:41.074847+00:00 app[web.1]: Diff block -  14108700
2022-01-30T18:05:41.413616+00:00 app[web.1]: time="2022-01-30T18:05:41Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:42.209820+00:00 app[web.1]: time="2022-01-30T18:05:42Z" level=info msg="Successful different block number"
2022-01-30T18:05:42.548585+00:00 app[web.1]: time="2022-01-30T18:05:42Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:43.683051+00:00 app[web.1]: time="2022-01-30T18:05:43Z" level=info msg="Current block number - 14108702"
2022-01-30T18:05:45.314831+00:00 app[web.1]: time="2022-01-30T18:05:45Z" level=info msg="Successful load block number - 14108702"
2022-01-30T18:05:45.314842+00:00 app[web.1]: time="2022-01-30T18:05:45Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:46.449800+00:00 app[web.1]: time="2022-01-30T18:05:46Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:47.586899+00:00 app[web.1]: time="2022-01-30T18:05:47Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:48.725578+00:00 app[web.1]: time="2022-01-30T18:05:48Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:49.860627+00:00 app[web.1]: time="2022-01-30T18:05:49Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:50.995576+00:00 app[web.1]: time="2022-01-30T18:05:50Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:52.130415+00:00 app[web.1]: time="2022-01-30T18:05:52Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:53.276133+00:00 app[web.1]: time="2022-01-30T18:05:53Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:54.411714+00:00 app[web.1]: time="2022-01-30T18:05:54Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:55.547469+00:00 app[web.1]: time="2022-01-30T18:05:55Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:56.683231+00:00 app[web.1]: time="2022-01-30T18:05:56Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:57.818178+00:00 app[web.1]: time="2022-01-30T18:05:57Z" level=info msg="Current block number - 14108703"
2022-01-30T18:05:58.953923+00:00 app[web.1]: time="2022-01-30T18:05:58Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:00.088678+00:00 app[web.1]: time="2022-01-30T18:06:00Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:01.222914+00:00 app[web.1]: time="2022-01-30T18:06:01Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:02.357420+00:00 app[web.1]: time="2022-01-30T18:06:02Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:03.492662+00:00 app[web.1]: time="2022-01-30T18:06:03Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:04.626875+00:00 app[web.1]: time="2022-01-30T18:06:04Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:05.761549+00:00 app[web.1]: time="2022-01-30T18:06:05Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:06.896704+00:00 app[web.1]: time="2022-01-30T18:06:06Z" level=info msg="Current block number - 14108703"
2022-01-30T18:06:08.528357+00:00 app[web.1]: time="2022-01-30T18:06:08Z" level=info msg="Successful load block number - 14108703"
2022-01-30T18:06:08.528383+00:00 app[web.1]: time="2022-01-30T18:06:08Z" level=info msg="Current block number - 14108704"
2022-01-30T18:06:09.662873+00:00 app[web.1]: time="2022-01-30T18:06:09Z" level=info msg="Current block number - 14108704"
2022-01-30T18:06:10.797399+00:00 app[web.1]: time="2022-01-30T18:06:10Z" level=info msg="Current block number - 14108704"
```  
