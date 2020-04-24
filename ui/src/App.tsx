
import React, { Component } from 'react';
import { Container, Divider, Button, Input, Form, List } from 'semantic-ui-react';
import 'semantic-ui-css/semantic.min.css?global';
import { Icon } from '@fluentui/react/lib/Icon';
import { initializeFileTypeIcons, getFileTypeIconProps } from '@uifabric/file-type-icons';
import './App.css';

initializeFileTypeIcons()

export default class App extends Component {

    render() {
        return (
            <div className='App'>
                <Container style={{width: '40%', paddingTop: 50}}>
                    <Form>
                        <Form.Group inline>
                            <Form.Field width='14'>
                                <Input size='mini' action fluid placeholder='Search...'>
                                    <input />
                                    <Button size='mini' icon='search' primary />
                                </Input>
                            </Form.Field>
                            <Form.Field width='2'>
                                <Button size='mini' circular icon='upload' primary />
                                <Button size='mini' circular icon='setting' primary />
                            </Form.Field>
                        </Form.Group>
                    </Form>
                    <Divider />
                    <List>
                        <List.Item>
                            <List.Icon verticalAlign='middle'>
                                <Icon {...getFileTypeIconProps({extension: 'docx', size: 16})} />
                            </List.Icon>
                            <List.Content>
                                <List.Header>test data file name</List.Header>
                                <List.Description>Uploaded 3 mins ago</List.Description>
                            </List.Content>
                        </List.Item>
                    </List>
                </Container>
            </div>
        )
    }
}