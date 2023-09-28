import React, { useEffect, useState } from 'react';
import { ColorPicker, Input } from 'antd';
import { doGetRequestAuth, doPostRequestAuth } from '../helper/RequestHelper';
import { myToastError, myToastSuccess } from '../helper/ToastHelper';

function handleSave(data, token) {
  const params = {skin: data.skin, gamename: data.gamename};
  doPostRequestAuth("user", params, token).then((e) => {
    if (e.status === 200) {
      myToastSuccess('Speichern erfolgreich');
    }}, error => {
      myToastError('Fehler beim speichern aufgetreten');
  })
}

function CharacterCreation(props) {
  const [value, setValue] = useState();

  useEffect(() => {
    doGetRequestAuth('user', props.token).then(
      res => {
        setValue(res.data)
      }
    )
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  return (
    <div>
      <div>
        <p>Username:</p>
        <Input value={value?.username} disabled></Input>
      </div>
      <br />
      <div>
        <p>Anzeigename:</p>
        <Input value={value?.gamename} onChange={(e)=>setValue({...value, gamename: e.target.value})}></Input>
      </div>
      <br />
      <div>
        <p>Farbe:</p>
        <ColorPicker
        value={value?.skin}
        onChangeComplete={(colorNew) => {
          setValue({...value, skin: colorNew.toHexString()});
        }}
      />
      </div>
      <br />
      <button type='primary' onClick={()=>handleSave(value, props.token)}>Speichern</button>
      
    </div>
  );
};
export default CharacterCreation;