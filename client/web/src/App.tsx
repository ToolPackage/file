
import React, { Component } from 'react';
import { Container, Divider, Button, Input, Form, List } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css?global';
import { Icon } from '@fluentui/react/lib/Icon';
import { initializeFileTypeIcons, getFileTypeIconProps } from '@uifabric/file-type-icons';
import Axios from './Axios';
import API from './API';
import './App.css';

initializeFileTypeIcons()

export default class App extends Component {

    componentDidMount() {
        Axios.get(API.getFileList)
            .then((rep) => {
                console.log(rep)
            })
        console.log(1)
    }

    render() {
        return (
            <div className='App'>
                <Container style={{width: '40%', paddingTop: 50}}>
                    <Form>
                        <Form.Group inline>
                            <Form.Field width='13'>
                                <Input size='mini' action fluid placeholder='Search...'>
                                    <input />
                                    <Button size='mini' icon='search' primary />
                                </Input>
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='upload' primary />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='download' primary />
                            </Form.Field>
                            <Form.Field>
                                <Button size='mini' circular icon='setting' primary />
                            </Form.Field>
                        </Form.Group>
                    </Form>
                    <Divider />
                    <List animated>
                        <List.Item style={{cursor: 'pointer'}}>
                            <List.Icon verticalAlign='middle'>
                                <Icon {...getFileTypeIconProps({extension: 'docx', size: 16})} />
                            </List.Icon>
                            <List.Content>
                                <List.Header>test data file name</List.Header>
                                <List.Description>
                                    <span>3kb</span>
                                    <span style={{marginLeft: 10}}>Uploaded 3 mins ago</span>
                                </List.Description>
                            </List.Content>
                        </List.Item>
                    </List>
                </Container>
            </div>
        )
    }
}