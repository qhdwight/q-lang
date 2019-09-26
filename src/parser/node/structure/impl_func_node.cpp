#include <iostream>

#include <boost/algorithm/string/split.hpp>
#include <boost/algorithm/string/trim_all.hpp>
#include <boost/algorithm/string/classification.hpp>

#include <util/terminal_color.hpp>
#include <parser/node/implementation/line_node.hpp>

#include "impl_func_node.hpp"

namespace ql::parser {
    ImplementFunctionNode::ImplementFunctionNode(std::string&& body, std::string_view const& innerBody, Tokens&& tokens,
                                                 AbstractNode::ParentRef const& parent)
            : ParseWithDescriptorNode(std::move(body), innerBody, std::move(tokens), parent) {
        registerVariableNode<LineNode>("default");
        registerVariableNode<IntegerNode>("int32");
    }

    void ImplementFunctionNode::parse() {
        // Split by tabs and spaces into tokens, which we use to find what type of node to create
        Tokens lines;
        boost::split(lines, m_InnerBody, boost::is_any_of(";"), boost::token_compress_on);
        // Remove first and last blank tokens if they exist
        boost::trim_all_if(lines, [](auto const& token) { return token.empty(); });
        for (auto line: lines) {
            std::cout << KCYN << line << RST << std::endl;
            Tokens tokens;
            boost::split(tokens, line, boost::is_any_of("\t "), boost::token_compress_on);
            boost::trim_all_if(tokens, [](auto const& token) { return token.empty(); });
            auto it = m_NameToVariableNode.find(tokens[0]);
            VariableNodeFactory& nodeFactoryFunc = it == m_NameToVariableNode.end() ? m_NameToVariableNode["default"] : it->second;
            m_Children.push_back(nodeFactoryFunc(std::move(line), std::move(tokens), weak_from_this()));
        }
    }

}
