import { Button, Table } from "antd";
import { useEffect, useState } from "react";
import { doGetRequest } from "../helper/RequestHelper";
import { useNavigate } from "react-router-dom";

function Lobby() {
  const navigate = useNavigate();
  const columns = [
    {
      title: 'Name',
      dataIndex: 'petName',
      key: 'petName',
    },
    {
      title: 'Address',
      dataIndex: 'address',
      key: 'address',
    },
    {
      title: 'Player Count',
      dataIndex: 'playerCount',
      key: 'playerCount',
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
        <Button onClick={()=>{
          navigate("/gameServer/"+record.address)
        }}>Beitreten</Button>
      ),
    },
  ];

  const [servers, setServers] = useState([]);

  useEffect(() => {
    doGetRequest('listServer').then(
      res => {
        setServers(
          res.data.map(row => ({
            key: row.id,
            id: row.id,
            petName : row.petName,
            address: row.address,
            playerCount: row.playerCount
          }))
        );
      }
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (    
    <Table columns={columns} dataSource={servers} />
  );
}

export default Lobby;
