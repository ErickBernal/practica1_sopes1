import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';
import axios from "axios";
import "bootstrap/dist/css/bootstrap.min.css";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCircleArrowRight, faEdit, faTrashAlt } from '@fortawesome/free-solid-svg-icons';
import { Modal, ModalBody, ModalFooter, ModalHeader } from 'reactstrap';


//const url_base = "http://192.168.122.210:8000"
const url_base = "http://localhost:8000"

const url_get = url_base+"/get"
const url_post = url_base+"/create"
const url_update = url_base+"/updateCar"
const url_delete = url_base+"/deleteCar"
const url_filterModelCar = url_base+"/filterCarModel"
const url_filterMarcaCar = url_base+"/filterCarMarca"
const url_filterColorCar = url_base+"/filterCarColor"



class App extends Component {
  state = {
    data: [],
    modalInsertar: false,
    modalEliminar: false,
    form: {
      placa: '',
      marca: '',
      modelo: '',
      serie: '',
      color: '',
      tipoModal: ''
    }
  }

  peticionGet = () => {
    axios.get(url_get).then(response => {
      console.log(response.data)
      if (response.data.modelo != "") {
        this.setState({ data: response.data });
      }
    }).catch(error => {
      console.log(error.message);
    })
  }

  peticionPost = async () => {
    await axios.post(url_post, this.state.form).then(response => {
      this.modalInsertar();
      this.peticionGet();
    }).catch(error => {
      console.log(error.message);
    })
  }

  seleccionarCarro = (carro) => {
    this.setState({
      tipoModal: 'actualizar',
      form: {
        placa: carro.placa,
        marca: carro.marca,
        modelo: carro.modelo,
        serie: carro.serie,
        color: carro.color
      }
    })
  }

  peticionPut = () => {
    axios.put(url_update, this.state.form).then(response => {
      this.modalInsertar();
      this.peticionGet();
    })
  }

  peticionDelete = () => {
    console.log(this.state.form)
    axios.post(url_delete, this.state.form).then(response => {
      this.setState({ modalEliminar: false });
      this.peticionGet();
    })
  }

  filtrarModelo = () => {
    console.log(this.state.form)
    axios.post(url_filterModelCar, this.state.form).then(response => {
      if (response.data.modelo != "") {
        this.setState({ data: response.data });
      }
    })
  }


  filtrarMarca = () => {
    console.log(this.state.form)
    axios.post(url_filterMarcaCar, this.state.form).then(response => {
      if (response.data.modelo != "") {
        this.setState({ data: response.data });
      }
    })
  }

  filtrarColor = () => {
    console.log(this.state.form)
    axios.post(url_filterColorCar, this.state.form).then(response => {
      if (response.data.modelo != "") {
        this.setState({ data: response.data });
      }
    })
  }



  modalInsertar = () => {
    this.setState({ modalInsertar: !this.state.modalInsertar });
  }



  handlechange = async e => {
    e.persist();
    await this.setState({
      form: {
        ...this.state.form,
        [e.target.name]: e.target.value
      }
    })
    console.log(this.state.form);
  }

  componentDidMount() {
    this.peticionGet();
  }



  render() {
    const { form } = this.state;
    return (
      <div className="App">
        <br />
        <tr>
          <button className="btn btn-success" onClick={() => { this.setState({ form: null, tipoModal: 'insertar' }); this.modalInsertar() }}>Crear Carro</button>
        </tr>
        <br />
        <tr>
          <button className="btn btn-secondary" onClick={() => { this.filtrarModelo() }}>Filtrar Modelo</button>
          {'--'}
          <button className="btn btn-warning" onClick={() => { this.filtrarColor() }}>Filtrar Color</button>
          {'--'}
          <button className="btn btn-primary" onClick={() => { this.filtrarMarca() }}>Filtrar Marca</button>
          <label htmlFor="modelo"></label>
          <input className="form-control" type="text" name="modelo" id="modelo" onChange={this.handlechange} value={form ? form.modelo : ''} />
        </tr>
        <br />

        <br />
        <table className="table">

          <thead>
            <tr>
              <th>PLACA</th>
              <th>MARCA</th>
              <th>MODELO</th>
              <th>SERIE</th>
              <th>COLOR</th>
            </tr>
          </thead>
          <tbody>
            {
              this.state.data.map(carro => {
                return (
                  <tr>
                    <td>{carro.placa}</td>
                    <td>{carro.marca}</td>
                    <td>{carro.modelo}</td>
                    <td>{carro.serie}</td>
                    <td>{carro.color}</td>
                    <td>
                      <button className="btn btn-primary" onClick={() => { this.seleccionarCarro(carro); this.modalInsertar() }}><FontAwesomeIcon icon={faEdit} /></button>
                      {"   "}
                      <button className="btn btn-danger" onClick={() => { this.seleccionarCarro(carro); this.setState({ modalEliminar: true }) }}><FontAwesomeIcon icon={faTrashAlt} /></button>
                    </td>
                  </tr>
                )
              })
            }
          </tbody>
        </table>


        <Modal isOpen={this.state.modalInsertar}>
          <ModalHeader style={{ display: 'block' }}>
            <span style={{ float: 'right' }} onClick={() => this.modalInsertar()}>x</span>
          </ModalHeader>
          <ModalBody>
            <div className="form-group">
              <label htmlFor="placa">Placa</label>
              <input className="form-control" type="text" name="placa" id="placa" onChange={this.handlechange} value={form ? form.placa : this.state.placa} />
              <br />
              <label htmlFor="marca">Marca</label>
              <input className="form-control" type="text" name="marca" id="marca" onChange={this.handlechange} value={form ? form.marca : ''} />
              <br />
              <label htmlFor="modelo">Modelo</label>
              <input className="form-control" type="text" name="modelo" id="modelo" onChange={this.handlechange} value={form ? form.modelo : ''} />
              <br />
              <label htmlFor="serie">Serie</label>
              <input className="form-control" type="text" name="serie" id="serie" onChange={this.handlechange} value={form ? form.serie : ''} />
              <label htmlFor="color">Color</label>
              <input className="form-control" type="text" name="color" id="color" onChange={this.handlechange} value={form ? form.color : ''} />
            </div>
          </ModalBody>

          <ModalFooter>
            {this.state.tipoModal == 'insertar' ?
              <button className="btn btn-success" onClick={() => this.peticionPost()}>
                Insertar
              </button> : <button className="btn btn-primary" onClick={() => this.peticionPut()}>
                Actualizar
              </button>
            }
            <button className="btn btn-danger" onClick={() => this.modalInsertar()}>Cancelar</button>
          </ModalFooter>
        </Modal>

        <Modal isOpen={this.state.modalEliminar}>
          <ModalBody >
            <tr> Esta seguro de eliminar el carro:</tr>
            <tr>Placa: </tr>
            <tr className="text-primary">"{form && form.placa}"</tr>
          </ModalBody>
          <ModalFooter>
            <button className="btn btn-danger" onClick={() => this.peticionDelete()}>Sí</button>
            <button className="btn btn-secundary" onClick={() => this.setState({ modalEliminar: false })}>No</button>
          </ModalFooter>
        </Modal>



      </div>
    );



  }
}

export default App;