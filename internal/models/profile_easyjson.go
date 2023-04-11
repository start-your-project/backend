// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson521a5691DecodeMainInternalModels(in *jlexer.Lexer, out *Recommend) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "profession":
			out.Profession = string(in.String())
		case "simularity":
			out.Simularity = float64(in.Float64())
		case "learned":
			if in.IsNull() {
				in.Skip()
				out.Learned = nil
			} else {
				in.Delim('[')
				if out.Learned == nil {
					if !in.IsDelim(']') {
						out.Learned = make([]string, 0, 4)
					} else {
						out.Learned = []string{}
					}
				} else {
					out.Learned = (out.Learned)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Learned = append(out.Learned, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "to_learn":
			if in.IsNull() {
				in.Skip()
				out.ToLearn = nil
			} else {
				in.Delim('[')
				if out.ToLearn == nil {
					if !in.IsDelim(']') {
						out.ToLearn = make([]string, 0, 4)
					} else {
						out.ToLearn = []string{}
					}
				} else {
					out.ToLearn = (out.ToLearn)[:0]
				}
				for !in.IsDelim(']') {
					var v2 string
					v2 = string(in.String())
					out.ToLearn = append(out.ToLearn, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeMainInternalModels(out *jwriter.Writer, in Recommend) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"profession\":"
		out.RawString(prefix[1:])
		out.String(string(in.Profession))
	}
	{
		const prefix string = ",\"simularity\":"
		out.RawString(prefix)
		out.Float64(float64(in.Simularity))
	}
	{
		const prefix string = ",\"learned\":"
		out.RawString(prefix)
		if in.Learned == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range in.Learned {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"to_learn\":"
		out.RawString(prefix)
		if in.ToLearn == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.ToLearn {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Recommend) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeMainInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Recommend) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeMainInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Recommend) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeMainInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Recommend) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeMainInternalModels(l, v)
}
func easyjson521a5691DecodeMainInternalModels1(in *jlexer.Lexer, out *ProfileUserDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "username":
			out.Name = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeMainInternalModels1(out *jwriter.Writer, in ProfileUserDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"avatar\":"
		out.RawString(prefix)
		out.String(string(in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ProfileUserDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeMainInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ProfileUserDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeMainInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ProfileUserDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeMainInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ProfileUserDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeMainInternalModels1(l, v)
}
func easyjson521a5691DecodeMainInternalModels2(in *jlexer.Lexer, out *LikeDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeMainInternalModels2(out *jwriter.Writer, in LikeDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LikeDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeMainInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LikeDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeMainInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LikeDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeMainInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LikeDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeMainInternalModels2(l, v)
}
func easyjson521a5691DecodeMainInternalModels3(in *jlexer.Lexer, out *Favorite) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "count_all":
			out.CountAll = int64(in.Int64())
		case "count_finished":
			out.CountFinished = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeMainInternalModels3(out *jwriter.Writer, in Favorite) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"count_all\":"
		out.RawString(prefix)
		out.Int64(int64(in.CountAll))
	}
	{
		const prefix string = ",\"count_finished\":"
		out.RawString(prefix)
		out.Int64(int64(in.CountFinished))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Favorite) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeMainInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Favorite) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeMainInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Favorite) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeMainInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Favorite) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeMainInternalModels3(l, v)
}
func easyjson521a5691DecodeMainInternalModels4(in *jlexer.Lexer, out *EmailUserDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeMainInternalModels4(out *jwriter.Writer, in EmailUserDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EmailUserDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeMainInternalModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EmailUserDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeMainInternalModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EmailUserDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeMainInternalModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EmailUserDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeMainInternalModels4(l, v)
}
func easyjson521a5691DecodeMainInternalModels5(in *jlexer.Lexer, out *EditProfileDTO) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "username":
			out.Name = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson521a5691EncodeMainInternalModels5(out *jwriter.Writer, in EditProfileDTO) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EditProfileDTO) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson521a5691EncodeMainInternalModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EditProfileDTO) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson521a5691EncodeMainInternalModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EditProfileDTO) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson521a5691DecodeMainInternalModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EditProfileDTO) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson521a5691DecodeMainInternalModels5(l, v)
}
